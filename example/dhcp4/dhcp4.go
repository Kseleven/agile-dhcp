package main

import (
	"flag"
	"github.com/Kseleven/agile-dhcp/dhcp4"
)

var (
	serverHost string
	hostName   string
	mac        string
)

func main() {
	flag.StringVar(&serverHost, "s", "255.255.255.255", "DHCP server IP")
	flag.StringVar(&hostName, "h", "", "client host name(option 12)")
	flag.StringVar(&mac, "m", "00:00:00:00:00:00", "client mac address(option 12)")

	flag.Parse()

	c, err := dhcp4.NewDHCPRequest("10.0.0.68", "joker", "00:00:00:00:00:01")
	if err != nil {
		panic(err)
	}
	c.WaitDone()
}
