package dhcp4

import (
	"fmt"
	"net"
	"strings"
)

type Conn struct {
	*net.UDPConn
	TransactionID      uint32
	SecondsElapsed     uint16
	DhcpServerHost     string
	CurrentMessageType MessageType
	HostName           string
	Mac                string
	MacByte            []byte
	doneChan           chan bool
}

func (c *Conn) Close() {
	if c.UDPConn != nil {
		c.UDPConn.Close()
	}
}

func NewDHCPRequest(serverIP string, hostName, mac string) (c *Conn, err error) {
	c = &Conn{
		DhcpServerHost: serverIP,
		SecondsElapsed: 0,
		TransactionID:  RandomTransactionID(),
		doneChan:       make(chan bool),
		Mac:            mac,
		HostName:       hostName,
	}

	hw, err := net.ParseMAC(c.Mac)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	c.MacByte = hw

	if serverIP == "" {
		serverIP = "255.255.255.255"
	}
	serverAddress := net.ParseIP(serverIP)
	if serverAddress == nil {
		return nil, fmt.Errorf("invalid server ip:%s", serverIP)
	}

	raddr := &net.UDPAddr{
		IP:   serverAddress,
		Port: 67,
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, fmt.Errorf("dial host %s failed:%s", raddr.IP, err.Error())
	}

	c.UDPConn = conn
	go c.listenUDP()

	return c, c.RequireAddress()
}

func (c *Conn) listenUDP() {
	laddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 68,
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Printf("listen udp failed:%s\n", err.Error())
		return
	}
	defer conn.Close()

	for {
		data := make([]byte, 576)
		length, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("read message failed:%s\n", err)
			if strings.Contains(err.Error(), "closed network connection") {
				return
			}
			continue
		}

		if ok := c.handlerResponse(rAddr, data[:length]); ok {
			c.done()
			return
		}
	}
}

func (c *Conn) WaitDone() {
	<-c.doneChan
	fmt.Println("wait done")
	c.Close()
}

func (c *Conn) done() {
	c.doneChan <- true
	fmt.Println("done")
}

func (c *Conn) RequireAddress() error {
	options := []OptionInter{
		GenOption51(7776000),
		GenOption57(1500),
		GenOption61(c.MacByte),
	}
	if c.HostName != "" {
		options = append(options, GenOption12(c.HostName))
	}

	m := GenDiscoverMessage(c.Mac, options...)
	m.TransactionID = c.TransactionID
	m.SecondsElapsed = c.SecondsElapsed
	c.CurrentMessageType = m.MessageType

	fmt.Printf("send message---->:\n%s\n", m.String())
	if _, err := c.Write(m.Encode()); err != nil {
		return err
	}

	return nil
}

func (c *Conn) handlerResponse(addr *net.UDPAddr, b []byte) bool {
	m := &Message{}
	m.Decode(b)

	if m.TransactionID != c.TransactionID {
		return false
	}
	fmt.Println("receive DHCP Message<----:", addr, len(b))
	fmt.Println(m.String())

	if m.MessageType == MessageTypeNak {
		return true
	}

	if c.CurrentMessageType == MessageTypeDiscover && m.MessageType == MessageTypeOffer {
		var options []OptionInter
		if c.HostName != "" {
			options = append(options, GenOption12(c.HostName))
		}
		requestMsg := GenRequestMessage(m, options...)
		c.SecondsElapsed = requestMsg.SecondsElapsed
		c.CurrentMessageType = requestMsg.MessageType
		fmt.Printf("send message---->:\n%s\n", m.String())
		if _, err := c.Write(requestMsg.Encode()); err != nil {
			fmt.Printf("write request message failed:%s\n", err.Error())
			return false
		}
	}

	if c.CurrentMessageType == MessageTypeRequest && m.MessageType == MessageTypeAck {
		return true
	}
	return false
}
