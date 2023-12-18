package tcpip

import (
	"errors"
	"log/slog"
	"net/netip"
	"sync"
	"time"

	"github.com/soypat/seqs/stacks"
	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
)

type sockets map[int]any // keyed by sockfd [1-n]

type Tcpip struct {
	stack *stacks.PortStack
	sockets
	socketsMu sync.RWMutex
}

func New(link netlink.Netlinker, logger *slog.Logger, MTU uint16) *Tcpip {
	t := Tcpip{}
	t.stack = stacks.NewPortStack(stacks.PortStackConfig{
		Link:            link,
//		Logger:          logger,
		MaxOpenPortsUDP: 1,
		MaxOpenPortsTCP: 1,
		MTU:             MTU,
	})
	t.sockets = make(sockets)
	return &t
}

func (t *Tcpip) GetHostByName(name string) (netip.Addr, error) {
	// Use ParseAddr to test if name is already in dotted decimal
	// ("10.0.0.1")
	addr, err := netip.ParseAddr(name)
	if err != nil {
		// Not in dotted-decimal
		// TODO implement
		return netip.Addr{}, netdev.ErrHostUnknown
	}
	return addr, nil
}

func (t *Tcpip) Addr() (netip.Addr, error) {
	return t.stack.Addr()
}

func (t *Tcpip) _newSockfd() int {
	var sockfd int

	// Find next available sockfd number, starting at 1
	for sockfd = 1;; sockfd++ {
		_, taken := t.sockets[sockfd]
		if !taken {
			break
		}
	}
	return sockfd
}

func (t *Tcpip) Socket(domain int, stype int, protocol int) (int, error) {

	println("Socket domain", domain, "stype", stype, "protocol", protocol)

	t.socketsMu.Lock()
	defer t.socketsMu.Unlock()

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

	sockfd := t._newSockfd()

	switch protocol {
	case netdev.IPPROTO_TCP:
		const socketBuf = 256
		sock, err := stacks.NewTCPSocket(t.stack, stacks.TCPSocketConfig{
			TxBufSize: socketBuf,
			RxBufSize: socketBuf,
		})
		if err != nil {
			return -1, err
		}
		t.sockets[sockfd] = sock
	default:
		return -1, netdev.ErrProtocolNotSupported
	}

	return sockfd, nil
}

func (t *Tcpip) Bind(sockfd int, ip netip.AddrPort) error {

	println("Bind sockfd", sockfd, "ip", ip.String())

	t.socketsMu.RLock()
	defer t.socketsMu.RUnlock()

	sock, found := t.sockets[sockfd]
	if !found {
		return netdev.ErrNoSocket
	}

	switch sock := sock.(type) {
	case *stacks.TCPSocket:
		return sock.Bind(ip)
	}

	return netdev.ErrNotSupported
}

func (t *Tcpip) Connect(sockfd int, host string, ip netip.AddrPort) error {

	println("Connect sockfd", sockfd, "host", host, "ip", ip.String())

	t.socketsMu.RLock()
	defer t.socketsMu.RUnlock()

	// TODO: for now fail host name connects
	if host != "" {
		return netdev.ErrNotSupported
	}

	sock, found := t.sockets[sockfd]
	if !found {
		return netdev.ErrNoSocket
	}

	switch sock := sock.(type) {
	case *stacks.TCPSocket:
		return sock.Connect(ip)
	}

	return netdev.ErrNotSupported
}

func (t *Tcpip) Listen(sockfd int, backlog int) error {

	println("Listen sockfd", sockfd, "backlog", backlog)

	t.socketsMu.RLock()
	defer t.socketsMu.RUnlock()

	sock, found := t.sockets[sockfd]
	if !found {
		return netdev.ErrNoSocket
	}

	switch sock := sock.(type) {
	case *stacks.TCPSocket:
		return sock.Listen(backlog)
	}

	return netdev.ErrNotSupported
}

func (t *Tcpip) Accept(sockfd int) (int, netip.AddrPort, error) {

	println("Accept sockfd", sockfd)

	t.socketsMu.Lock()
	defer t.socketsMu.Unlock()

	sock, found := t.sockets[sockfd]
	if !found {
		return -1, netip.AddrPort{}, netdev.ErrNoSocket
	}

	newSockfd := t._newSockfd()

	switch sock := sock.(type) {
	case *stacks.TCPSocket:
		newSock, raddr, err := sock.Accept()
		if err != nil {
			return -1, netip.AddrPort{}, err
		}
		t.sockets[newSockfd] = newSock
		println("Accept sockfd", sockfd, "--> New sockfd", newSockfd)
		return newSockfd, raddr, nil
	}

	return -1, netip.AddrPort{}, netdev.ErrNotSupported
}

func (t *Tcpip) Send(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {
	t.socketsMu.RLock()
	defer t.socketsMu.RUnlock()

	return 0, errors.New("Send not implemented")
}

func (t *Tcpip) Recv(sockfd int, buf []byte, flags int, deadline time.Time) (int, error) {
	return 0, errors.New("Recv not implemented")
}

func (t *Tcpip) Close(sockfd int) error {
	t.socketsMu.Lock()
	defer t.socketsMu.Unlock()

	return errors.New("Close not implemented")
}

func (t *Tcpip) SetSockOpt(sockfd int, level int, opt int, value interface{}) error {
	t.socketsMu.RLock()
	defer t.socketsMu.RUnlock()

	return errors.New("SetSockOpt not implemented")
}
