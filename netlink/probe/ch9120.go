//go:build rp2040_eth

package probe

import (
	"machine"

	"tinygo.org/x/drivers/ch9120"
	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
)

func Probe() (netlink.Netlinker, netdev.Netdever) {

	ch9120 := ch9120.NewDevice(&ch9120.Config{
		Uart:    machine.UART1,
		Tx:      machine.GP20,
		Rx:      machine.GP21,
		Cfg:     machine.GP18,
		Rst:     machine.GP19,
		RunBaud: 115200,
		//RunBaud: 921600,
	})

	netdev.UseNetdev(ch9120)

	return ch9120, ch9120
}
