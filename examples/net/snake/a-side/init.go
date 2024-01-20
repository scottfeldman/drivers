//go:build pyportal || nano_rp2040 || metro_m4_airlift || arduino_mkrwifi1010 || matrixportal_m4 || wioterminal

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
