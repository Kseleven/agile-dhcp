package main

import (
	"flag"
	"github.com/Kseleven/agile-dhcp/dhcp4"
)

var (
	serverHost string
	hostName   string
	relay      string
	mac        string
	count      int
	decline    string
	release    string
)

func main() {
	flag.StringVar(&serverHost, "s", "255.255.255.255", "DHCP server IP")
	flag.StringVar(&hostName, "h", "", "client host name(option 12)")
	flag.StringVar(&relay, "g", "", "relay ip")
	flag.StringVar(&decline, "d", "", "decline address")
	flag.StringVar(&release, "r", "", "release address")
	flag.StringVar(&mac, "m", "00:00:00:00:00:00", "client mac address(option 12)")
	flag.IntVar(&count, "c", 1, "numbers client")
	flag.Parse()

	if decline != "" {
		c, err := dhcp4.NewDHCPRequest(serverHost, relay, hostName, mac)
		if err != nil {
			panic(err)
		}
		if err := c.Decline(decline); err != nil {
			panic(err)
		}
		c.WaitDone()
		return
	}

	if release != "" {
		c, err := dhcp4.NewDHCPRequest(serverHost, relay, hostName, mac)
		if err != nil {
			panic(err)
		}
		if err := c.Release(release); err != nil {
			panic(err)
		}
		c.WaitDone()
		return
	}

	for i := 0; i < count; i++ {
		c, err := dhcp4.NewDHCPRequest(serverHost, relay, hostName, mac)
		if err != nil {
			panic(err)
		}
		if err := c.Discovery(); err != nil {
			panic(err)
		}
		c.WaitDone()
	}
}
