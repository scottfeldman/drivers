// Package ch9120 implements TCP/UDP wired communication over serial UART using
// CH9120 network chip.  The chip has an embedded TCP/UDP stack which supports
// only a single socket at a time in one of four configurations: TCP/UDP
// client/server.
//
// Datasheet: https://files.waveshare.com/upload/d/d3/CH9120DS1_EN.pdf
// Instruction Set: https://files.waveshare.com/upload/e/e1/CH9120_Serial_Commands_Set.pdf

package ch9120 // import "tinygo.org/x/drivers/ch9120"

import (
	"fmt"
	"machine"
	"net"
	"net/netip"
	"time"

	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
)

const (
	cmdVersion          = 0x01
	cmdReset            = 0x02
	cmdGetStatus        = 0x03
	cmdSaveEeprom       = 0x0d
	cmdExecCfg          = 0x0e
	cmdSetMode          = 0x10
	cmdSetSrcIp         = 0x11
	cmdSetSubnet        = 0x12
	cmdSetGateway       = 0x13
	cmdSetSrcPort       = 0x14
	cmdSetDstIp         = 0x15
	cmdSetDstPort       = 0x16
	cmdSetSrcPortRandom = 0x17
	cmdSetBaud          = 0x21
	cmdSetTimeout       = 0x23
	cmdSetDisconnect    = 0x24
	cmdSetRxPktLength   = 0x25
	cmdSetSerialClear   = 0x26
	cmdSetDhcp          = 0x33
	cmdExit             = 0x5e
	cmdGetIp            = 0x61
	cmdGetSubnet        = 0x62
	cmdGetGateway       = 0x63
	cmdGetDstIp         = 0x65
	cmdGetTimeout       = 0x73
	cmdGetDisconnect    = 0x74
	cmdGetRxPktLength   = 0x75
	cmdGetSerialClear   = 0x76
	cmdGetMac           = 0x81
)

const (
	maxSendSize = 512
)

var (
	tcpServer = []byte{0x00}
	tcpClient = []byte{0x01}
	udpServer = []byte{0x02}
	udpClient = []byte{0x03}
)

var (
	noParams       = []byte{}
	disconnectNo   = []byte{0x00}
	disconnectYes  = []byte{0x01}
	serialClearNo  = []byte{0x00}
	serialClearYes = []byte{0x01}
	noTimeout      = []byte{0x00, 0x00, 0x00, 0x00}
	randomSrcPort  = []byte{0x01}
	dhcpOff        = []byte{0x00}
	dhcpOn         = []byte{0x01}
)

type Config struct {
	Uart    *machine.UART
	Tx      machine.Pin
	Rx      machine.Pin
	Cfg     machine.Pin
	Rst     machine.Pin
	RunBaud uint32
}

type socket struct {
	inUse    bool
	protocol int
	laddr    netip.AddrPort
}

type Device struct {
	uart    *machine.UART
	cfg     machine.Pin
	rst     machine.Pin
	cfgBaud uint32
	runBaud uint32
	mac     net.HardwareAddr
	ip      netip.Addr
	socket  socket
	cmdBuf  [32]byte
}

func NewDevice(cfg *Config) *Device {
	d := Device{
		uart:    cfg.Uart,
		cfg:     cfg.Cfg,
		rst:     cfg.Rst,
		cfgBaud: 9600,
		runBaud: cfg.RunBaud,
	}
	d.uart.Configure(machine.UARTConfig{TX: cfg.Tx, RX: cfg.Rx})
	d.cfg.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.rst.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return &d
}

func (d *Device) NetConnect(params *netlink.ConnectParams) error {

	d.reset()

	d.setBaud(d.cfgBaud)
	d.cfgBegin()

	// On network cable disconnect, disconnect network connection
	d.cmd(cmdSetDisconnect, disconnectYes)
	// On network (re)connection, clear serial port
	d.cmd(cmdSetSerialClear, serialClearYes)
	// No timeout on send; send right away
	d.cmd(cmdSetTimeout, noTimeout)
	//d.cmd(cmdGetTimeout, noParams)

	mac, _ := d.cmd(cmdGetMac, noParams)
	d.mac = net.HardwareAddr(mac)
	println("MAC:        ", d.mac.String())

	d.cmd(cmdSetDhcp, dhcpOff)

	// get IP, subnet, gateway
	for {
		var ok bool
		addr, _ := d.cmd(cmdGetIp, noParams)
		d.ip, ok = netip.AddrFromSlice(addr)
		if ok {
			break
		}
	}
	addr, _ := d.cmd(cmdGetSubnet, noParams)
	subnet, _ := netip.AddrFromSlice(addr)
	addr, _ = d.cmd(cmdGetGateway, noParams)
	gateway, _ := netip.AddrFromSlice(addr)

	println("IP:         ", d.ip.String())
	println("Subnet:     ", subnet.String())
	println("Gateway:    ", gateway.String())

	d.save()
	d.cfgEnd()

	return nil
}

func (d *Device) NetDisconnect() {
}

func (d *Device) NetNotify(cb func(netlink.Event)) {
	// Not supported
}

func (d *Device) GetHostByName(name string) (netip.Addr, error) {
	// If host name is dotted-decimal notation, just return as netip.Addr
	ip, err := netip.ParseAddr(name)
	if err == nil {
		return ip, nil
	}

	// Set dst IP address to host name and then read back resolved
	// dst IP address
	d.reset()
	d.setBaud(d.cfgBaud)
	d.cfgBegin()
	d.cmd(cmdSetDstIp, []byte(name))
	addr, _ := d.cmd(cmdGetDstIp, noParams)
	ip, _ = netip.AddrFromSlice(addr)
	d.cfgEnd()

	return ip, nil
}

func (d *Device) GetHardwareAddr() (net.HardwareAddr, error) {
	return d.mac, nil
}

func (d *Device) Addr() (netip.Addr, error) {
	return d.ip, nil
}

func (d *Device) Socket(domain int, stype int, protocol int) (int, error) {

	switch domain {
	case netdev.AF_INET:
	default:
		return -1, netdev.ErrFamilyNotSupported
	}

	switch {
	case protocol == netdev.IPPROTO_TCP && stype == netdev.SOCK_STREAM:
	case protocol == netdev.IPPROTO_UDP && stype == netdev.SOCK_DGRAM:
	default:
		return -1, netdev.ErrProtocolNotSupported
	}

	// Only supporting single connection mode, so only one socket at a time
	if d.socket.inUse {
		return -1, netdev.ErrNoMoreSockets
	}
	d.socket.inUse = true
	d.socket.protocol = protocol

	return 1, nil
}

func (d *Device) Bind(sockfd int, ip netip.AddrPort) error {
	d.socket.laddr = ip
	return nil
}

func (d *Device) tcpConnect(ip netip.AddrPort) {
	d.reset()
	d.setBaud(d.cfgBaud)
	d.cfgBegin()
	// start TCP client
	d.cmd(cmdSetMode, tcpClient)
	// use random (ephemeral) local src port
	d.cmd(cmdSetSrcPortRandom, randomSrcPort)
	// set dst ip:port
	raddr := ip.Addr().AsSlice()
	rport := ip.Port()
	d.cmd(cmdSetDstIp, raddr)
	d.cmd(cmdSetDstPort, port(rport))
	// set rx/tx baudrate
	d.cmd(cmdSetBaud, baud(d.runBaud))
	d.save()
	d.cfgEnd()
	d.setBaud(d.runBaud)
	// ready for tx/rx
}

func (d *Device) Connect(sockfd int, host string, ip netip.AddrPort) error {

	switch d.socket.protocol {
	case netdev.IPPROTO_TCP:
		d.tcpConnect(ip)
	}

	return nil
}

func (d *Device) Listen(sockfd int, backlog int) error {
	return netdev.ErrNotSupported
}

func (d *Device) Accept(sockfd int, ip netip.AddrPort) (int, error) {
	return -1, netdev.ErrNotSupported
}

func (d *Device) Send(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {

	// Break large bufs into chunks so we don't overrun the hw queue

	chunkSize := maxSendSize
	for i := 0; i < len(buf); i += chunkSize {
		end := i + chunkSize
		if end > len(buf) {
			end = len(buf)
		}
		_, err := d.uart.Write(buf[i:end])
		if err != nil {
			return -1, err
		}
	}

	return len(buf), nil
}

func (d *Device) Recv(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {
	n, err := d.uart.Read(buf)
	println("recv", n, err)
	if n > 0 {
		println(string(buf[:n]))
	}
	time.Sleep(time.Second)
	return n, err
}

func (d *Device) Close(sockfd int) error {
	time.Sleep(5 * time.Millisecond)
	d.reset()
	d.socket.inUse = false
	return nil
}

func (d *Device) SetSockOpt(sockfd int, level int, opt int, value interface{}) error {
	return netdev.ErrNotSupported
}

// save config to eeprom
func (d *Device) save() {
	d.cmd(cmdSaveEeprom, noParams)
	d.cmd(cmdExecCfg, noParams)
	d.cmd(cmdExit, noParams)
	time.Sleep(100 * time.Millisecond)
}

// Drain UART receive buffer
func (d *Device) drain() {
	for d.uart.Buffered() > 0 {
		d.uart.ReadByte()
	}
}

// Reset hard
func (d *Device) reset() {
	d.rst.Low()
	time.Sleep(500 * time.Millisecond)
	d.rst.High()
	time.Sleep(500 * time.Millisecond)
}

func (d *Device) setBaud(rate uint32) {
	d.uart.SetBaudRate(rate)
}

// Begin serial port configuration mode in low-speed cfgBaud
func (d *Device) cfgBegin() {
	d.cfg.Low()
	time.Sleep(100 * time.Millisecond)
	d.drain()
}

// End serial port configuration mode and start I/O mode at full runBaud speed
func (d *Device) cfgEnd() {
	d.cfg.High()
	time.Sleep(100 * time.Millisecond)
	d.drain()
}

// Run command with optional params and return response
func (d *Device) cmd(cmdCode uint8, params []byte) ([]byte, error) {

	// The format of the command code sent by CH9121 is:
	// "0x57 0xab commandcode parameter (optional)"

	send := []byte{0x57, 0xab, cmdCode}
	send = append(send, params...)

	n, err := d.uart.Write(send)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("UART write failed: %w", err)
	}
	//println("write", n, hex.EncodeToString(send[0:2]), hex.EncodeToString(send[2:3]), hex.EncodeToString(send[3:n]))

	n, err = d.uart.Read(d.cmdBuf[:])
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("UART read failed: %w", err)
	}
	//println("read", n)
	//println(hex.Dump(d.cmdBuf[:n]))

	return d.cmdBuf[:n], nil
}

func ip(a, b, c, d byte) []byte {
	return []byte{a, b, c, d}
}

func port(p uint16) []byte {
	return []byte{byte(p), byte(p >> 8)}
}

func baud(b uint32) []byte {
	return []byte{byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24)}
}
