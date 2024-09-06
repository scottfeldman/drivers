//go:build ninafw && !arduino_mkrwifi1010

package probe

import (
	"crypto/rand"
	"machine"

	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/wifinina"
)

func init() {
	rand.Reader = &reader{}
}

type reader struct{}

func (r *reader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return
	}
	var randomByte uint32
	for i := range b {
		if i%4 == 0 {
			randomByte, err = machine.GetRNG()
			if err != nil {
				return n, err
			}
		} else {
			randomByte >>= 8
		}
		b[i] = byte(randomByte)
	}
	return len(b), nil
}

func Probe() (netlink.Netlinker, netdev.Systemer) {

	cfg := wifinina.Config{
		// Configure SPI for 8Mhz, Mode 0, MSB First
		Spi:  machine.NINA_SPI,
		Freq: 8 * 1e6,
		Sdo:  machine.NINA_SDO,
		Sdi:  machine.NINA_SDI,
		Sck:  machine.NINA_SCK,
		// Device pins
		Cs:     machine.NINA_CS,
		Ack:    machine.NINA_ACK,
		Gpio0:  machine.NINA_GPIO0,
		Resetn: machine.NINA_RESETN,
	}

	nina := wifinina.New(&cfg)
	netdev.UseSystem(nina)

	return nina, nina
}
