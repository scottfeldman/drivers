//go:build pico

package probe

import (
	"log/slog"
	"machine"

	"github.com/soypat/cyw43439"
	"tinygo.org/x/drivers/netdev"
	"tinygo.org/x/drivers/netdev/tcpip"
	"tinygo.org/x/drivers/netlink"
)

const MTU = cyw43439.MTU // CY43439 permits 2030 bytes of ethernet frame.

func Probe() (netlink.Netlinker, netdev.Netdever) {

	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Go lower (Debug-1) to see more verbosity on wifi device.
	}))

	link := cyw43439.NewPicoWDevice(logger)
	stack := tcpip.New(link, logger, MTU)
	netdev.UseNetdev(stack)

	return link, stack
}
