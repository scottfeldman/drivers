//go:build ninafw || wioterminal || elecrow_rp2040 || elecrow_rp2350

package main

import (
	"log"
	"time"

	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"
)

var (
	ssid string
	pass string
)

func init() {
	time.Sleep(2 * time.Second)

	link, _ := probe.Probe()

	err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       ssid,
		Passphrase: pass,
	})
	if err != nil {
		log.Fatal(err)
	}
}
