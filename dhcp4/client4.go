package dhcp4

import (
	"fmt"
	"net"
	"time"
)

const MaxRetryNum = 1

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
	retry              int
	relay              []byte
	ifnname            *net.Interface
}

func (c *Conn) Close() {
	if c.UDPConn != nil {
		c.UDPConn.Close()
	}
}

func (c *Conn) isRelay() bool {
	return !(c.relay[0] == 0 && c.relay[1] == 0 && c.relay[2] == 0 && c.relay[3] == 0)
}

func NewDHCPRequest(serverIP string, relay string, hostName, mac string) (c *Conn, err error) {
	c = &Conn{
		DhcpServerHost: serverIP,
		SecondsElapsed: 0,
		TransactionID:  RandomTransactionID(),
		doneChan:       make(chan bool),
		Mac:            mac,
		HostName:       hostName,
		relay:          make([]byte, 4, 4),
	}

	hw, err := net.ParseMAC(c.Mac)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	c.MacByte = hw

	if c.ifnname, err = net.InterfaceByName("en0"); err != nil {
		return nil, err
	}

	if relay != "" {
		if addr := net.ParseIP(relay); addr == nil || addr.To4() == nil {
			return nil, fmt.Errorf("invalid relay ip")
		} else {
			c.relay = addr.To4()
		}
	}

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
	return c, nil
}

func (c *Conn) listenUDP() {
	laddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 68,
	}
	if c.isRelay() {
		laddr.Port = 67
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Printf("listen udp failed:%s\n", err.Error())
		return
	}
	defer conn.Close()

	now := time.Now()
	conn.SetReadDeadline(now.Add(time.Second * 3))
	for {
		data := make([]byte, 576)
		length, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("read message failed:%s\n", err)
			if op, ok := err.(*net.OpError); ok && (op.Timeout() || op.Temporary()) {
				c.done()
				return
			}
			continue
		}

		if ok := c.handlerResponse(rAddr, data[:length]); ok {
			c.done()
			return
		}

		now = time.Now()
		conn.SetReadDeadline(now.Add(time.Second * 3))
	}
}

func (c *Conn) WaitDone() {
	<-c.doneChan
	c.Close()
}

func (c *Conn) done() {
	c.doneChan <- true
	fmt.Println("done")
}

func (c *Conn) Discovery() error {
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
	m.RelayAgentIP = c.relay

	fmt.Printf("send message---->:\n%s\n", m.String())
	if _, err := c.Write(m.Encode()); err != nil {
		return err
	}

	return nil
}

func (c *Conn) Decline(declineIP string) error {
	server := net.ParseIP(c.DhcpServerHost)
	options := []OptionInter{
		GenOption54(server.To4()),
		GenOption57(1500),
		GenOption61(c.MacByte),
	}
	if c.HostName != "" {
		options = append(options, GenOption12(c.HostName))
	}

	declineIp := net.ParseIP(declineIP)
	m := GenDeclineMessage(c.Mac, declineIp.To4(), options...)
	m.TransactionID = c.TransactionID
	m.SecondsElapsed = c.SecondsElapsed
	c.CurrentMessageType = m.MessageType
	m.RelayAgentIP = c.relay

	fmt.Printf("send message---->:\n%s\n", m.String())
	if _, err := c.Write(m.Encode()); err != nil {
		return err
	}

	return nil
}

func (c *Conn) Release(release string) error {
	server := net.ParseIP(c.DhcpServerHost)
	options := []OptionInter{
		GenOption54(server.To4()),
		GenOption57(1500),
		GenOption61(c.MacByte),
	}
	if c.HostName != "" {
		options = append(options, GenOption12(c.HostName))
	}

	releaseIP := net.ParseIP(release)
	m := GenReleaseMessage(c.Mac, releaseIP.To4(), options...)
	m.TransactionID = c.TransactionID
	m.SecondsElapsed = c.SecondsElapsed
	c.CurrentMessageType = m.MessageType
	m.RelayAgentIP = c.relay

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
		c.retry++
		if c.CurrentMessageType == MessageTypeDiscover && c.retry < MaxRetryNum {
			if err := c.Discovery(); err != nil {
				fmt.Printf("write request message failed:%s\n", err.Error())
				return false
			}
			return false
		}

		if addr.String() != c.DhcpServerHost {
			return false
		}
		return true
	}

	c.retry = 0
	if c.CurrentMessageType == MessageTypeDiscover && m.MessageType == MessageTypeOffer {
		options := []OptionInter{
			GenOption57(1500),
			GenOption51(7776000),
		}

		if c.HostName != "" {
			options = append(options, GenOption12(c.HostName))
		}
		requestMsg := GenRequestMessage(m, options...)
		c.SecondsElapsed = requestMsg.SecondsElapsed
		c.CurrentMessageType = requestMsg.MessageType
		requestMsg.RelayAgentIP = c.relay

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
