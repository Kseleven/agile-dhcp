package main

import (
	"flag"
	"github.com/Kseleven/agile-dhcp/dhcp4"
)

var (
	serverHost string
	hostName   string
	mac        string
	count      int
)

func main() {
	flag.StringVar(&serverHost, "s", "255.255.255.255", "DHCP server IP")
	flag.StringVar(&hostName, "h", "", "client host name(option 12)")
	flag.StringVar(&mac, "m", "00:00:00:00:00:00", "client mac address(option 12)")
	flag.IntVar(&count, "c", 1, "numbers client")

	flag.Parse()

	for i := 0; i < count; i++ {
		c, err := dhcp4.NewDHCPRequest(serverHost, hostName, mac)
		if err != nil {
			panic(err)
		}
		c.WaitDone()
	}
}
