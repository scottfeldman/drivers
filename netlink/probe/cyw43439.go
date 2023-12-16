//go:build pico

package probe

import (
	"log/slog"
	"machine"

	"github.com/soypat/cyw43439"
	"github.com/soypat/seqs/eth"
	"github.com/soypat/seqs/stacks"
	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netlink"
)

const MTU = cyw43439.MTU // CY43439 permits 2030 bytes of ethernet frame.

func Probe() (netlink.Netlinker, netdev.Netdever) {

	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Go lower (Debug-1) to see more verbosity on wifi device.
	}))

	link := cyw43439.NewPicoWDevice(logger)

	dev := stacks.NewPortStack(stacks.PortStackConfig{
		MaxOpenPortsUDP: 1,
		MaxOpenPortsTCP: 1,
		GlobalHandler: func(ehdr *eth.EthernetHeader, ethPayload []byte) error {
			//lastRx = time.Now()
			return nil
		},
		MTU:    MTU,
		Logger: logger,
		Link:   link,
	})
	netdev.UseNetdev(dev)

	return link, dev
}
