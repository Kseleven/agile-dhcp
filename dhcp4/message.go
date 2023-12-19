package dhcp4

import (
	"bytes"
	"encoding/hex"
	"net"
)

type Message struct {
	OpCode         uint8          `json:"op"`        //op(1 octet):Message op code / message type 1 = BOOTREQUEST, 2 = BOOTREPLY
	HardwareType   uint8          `json:"htype"`     //htype(1 octet): Hardware address type, see ARP section in "Assigned Numbers" RFC; e.g., ’1’ = 10mb ethernet
	HardwareLength uint8          `json:"hlen"`      //hlen(1 octet): Hardware address length(e.g.  ’6’ for 10mb ethernet)
	Hops           uint8          `json:"hops"`      //hops(1 octet):Client sets to zero, optionally used by relay agents when booting via a relay agent.
	TransactionID  uint32         `json:"xid"`       //xid(4 octets):Transaction ID, a random number chosen by the client, used by the client and server to associate messages and responses between a client and a server
	SecondsElapsed uint16         `json:"secs"`      //secs(2 octets):Filled in by client, seconds elapsed since client began address acquisition or renewal process.
	Flags          uint16         `json:"flags"`     //flags(2 octets):Bootp Flags
	ClientIP       []byte         `json:"ciaddr"`    //ciaddr(4 octets):Client IP address; only filled in if client is in BOUND, RENEW or REBINDING state and can respond to ARP requests
	YourIP         []byte         `json:"yiaddr"`    //yiaddr(4 octets):’your’ (client) IP address
	NextServerIP   []byte         `json:"siaddr"`    //siaddr(4 octets):IP address of next server to use in bootstrap;returned in DHCPOFFER, DHCPACK by server.
	RelayAgentIP   []byte         `json:"giaddr"`    //giaddr(4 octets):Relay agent IP address, used in booting via a relay agent
	ClientMAC      ClientHardware `json:"chaddr"`    //chaddr(16 octets):Client hardware address(6 octets)+Client hardware address padding(10 octets)
	ServerHostName []byte         `json:"sname"`     //sname(64 octets):Optional server host name, null terminated string
	BootFile       []byte         `json:"file"`      //file(128 octets):Boot file name, null terminated string; "generic" name or null in DHCPDISCOVER, fully qualified directory-path name in DHCPOFFER
	MagicCookie    []byte         `json:"magicDhcp"` //magicDhcp(4 octets):fixed value[63 92 53 63]
	Options        []OptionInter  `json:"options"`   //options(var):Optional parameters field
	MessageType    MessageType
}

var MagicCookie = []byte{0x63, 0x82, 0x53, 0x63}

const MinRequestLength = 300

type ClientHardware struct {
	HardwareAddress        []byte `json:"chaddr"`        //Client hardware address(6 octets)
	HardwareAddressPadding []byte `json:"chaddrpadding"` //Client hardware address padding(10 octets)
}

func GenClientHardware(mac string) (ClientHardware, error) {
	m, err := net.ParseMAC(mac)
	if err != nil {
		return ClientHardware{}, err
	}

	return ClientHardware{HardwareAddress: m, HardwareAddressPadding: make([]byte, 10, 10)}, nil
}

func (c ClientHardware) Encode() []byte {
	var buf bytes.Buffer
	buf.Write(c.HardwareAddress)
	buf.Write(c.HardwareAddressPadding)
	return buf.Bytes()
}

func (m *Message) getOption(code uint8) OptionInter {
	for _, option := range m.Options {
		if option.GetCode() == code {
			return option
		}
	}

	return nil
}

func GenDiscoverMessage(mac string, options ...OptionInter) *Message {
	m := &Message{}
	m.OpCode = 1
	m.HardwareType = 1
	m.HardwareLength = 6
	m.Hops = 0
	m.TransactionID = 0
	m.SecondsElapsed = 0
	m.Flags = 0
	m.ClientIP = make([]byte, 4, 4)
	m.YourIP = make([]byte, 4, 4)
	m.NextServerIP = make([]byte, 4, 4)
	m.RelayAgentIP = make([]byte, 4, 4)
	m.ClientMAC, _ = GenClientHardware(mac)
	m.ServerHostName = make([]byte, 64, 64)
	m.BootFile = make([]byte, 128, 128)
	m.MagicCookie = MagicCookie
	m.Options = []OptionInter{GenOption53(MessageTypeDiscover), GenOption55()}
	for _, option := range options {
		m.Options = append(m.Options, option)
	}
	m.Options = append(m.Options, GenOption255())
	m.MessageType = MessageTypeDiscover
	return m
}

func GenDeclineMessage(mac string, declineIP []byte, options ...OptionInter) *Message {
	m := &Message{}
	m.OpCode = 1
	m.HardwareType = 1
	m.HardwareLength = 6
	m.Hops = 0
	m.TransactionID = 0
	m.SecondsElapsed = 0
	m.Flags = 0
	m.ClientIP = make([]byte, 4, 4)
	m.YourIP = make([]byte, 4, 4)
	m.NextServerIP = make([]byte, 4, 4)
	m.RelayAgentIP = make([]byte, 4, 4)
	m.ClientMAC, _ = GenClientHardware(mac)
	m.ServerHostName = make([]byte, 64, 64)
	m.BootFile = make([]byte, 128, 128)
	m.MagicCookie = MagicCookie
	m.Options = []OptionInter{GenOption53(MessageTypeDecline), GenOption55(), GenOption50(declineIP)}
	for _, option := range options {
		m.Options = append(m.Options, option)
	}
	m.Options = append(m.Options, GenOption255())
	m.MessageType = MessageTypeDecline
	return m
}

func GenReleaseMessage(mac string, releaseIP []byte, options ...OptionInter) *Message {
	m := &Message{}
	m.OpCode = 1
	m.HardwareType = 1
	m.HardwareLength = 6
	m.Hops = 0
	m.TransactionID = 0
	m.SecondsElapsed = 0
	m.Flags = 0
	m.ClientIP = releaseIP
	m.YourIP = make([]byte, 4, 4)
	m.NextServerIP = make([]byte, 4, 4)
	m.RelayAgentIP = make([]byte, 4, 4)
	m.ClientMAC, _ = GenClientHardware(mac)
	m.ServerHostName = make([]byte, 64, 64)
	m.BootFile = make([]byte, 128, 128)
	m.MagicCookie = MagicCookie
	m.Options = []OptionInter{GenOption53(MessageTypeRelease), GenOption55()}
	for _, option := range options {
		m.Options = append(m.Options, option)
	}
	m.Options = append(m.Options, GenOption255())
	m.MessageType = MessageTypeRelease
	return m
}

func GenRequestMessage(offer *Message, options ...OptionInter) *Message {
	m := &Message{}
	m.OpCode = 1
	m.HardwareType = 1
	m.HardwareLength = 6
	m.Hops = 0
	m.TransactionID = offer.TransactionID
	m.SecondsElapsed = offer.SecondsElapsed + 1
	m.Flags = 0
	m.ClientIP = make([]byte, 4, 4)
	m.YourIP = make([]byte, 4, 4)
	m.NextServerIP = make([]byte, 4, 4)
	m.RelayAgentIP = make([]byte, 4, 4)
	m.ClientMAC = offer.ClientMAC
	m.ServerHostName = make([]byte, 64, 64)
	m.BootFile = make([]byte, 128, 128)
	m.MagicCookie = MagicCookie
	m.Options = []OptionInter{GenOption53(MessageTypeRequest), GenOption55(), GenOption50(offer.YourIP)}
	for _, option := range options {
		m.Options = append(m.Options, option)
	}
	if option54 := offer.getOption(54); option54 != nil {
		m.Options = append(m.Options, option54)
	}
	if option61 := offer.getOption(61); option61 != nil {
		m.Options = append(m.Options, option61)
	}
	m.Options = append(m.Options, GenOption255())
	m.MessageType = MessageTypeRequest
	return m
}

func (m *Message) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(m.OpCode)
	buf.WriteByte(m.HardwareType)
	buf.WriteByte(m.HardwareLength)
	buf.WriteByte(m.Hops)
	buf.Write(Uint32ToBytes(m.TransactionID))
	buf.Write(Uint16ToBytes(m.SecondsElapsed))
	buf.Write(Uint16ToBytes(m.Flags))
	buf.Write(m.ClientIP)
	buf.Write(m.YourIP)
	buf.Write(m.NextServerIP)
	buf.Write(m.RelayAgentIP)
	buf.Write(m.ClientMAC.Encode())
	buf.Write(m.ServerHostName)
	buf.Write(m.BootFile)
	buf.Write(m.MagicCookie)
	for _, option := range m.Options {
		buf.Write(option.Encode())
	}

	diff := MinRequestLength - buf.Len()
	for diff > 0 {
		buf.WriteByte(0)
		diff--
	}
	return buf.Bytes()
}

func (m *Message) Decode(data []byte) {
	var buf bytes.Buffer
	buf.Grow(len(data))
	buf.Write(data)

	m.OpCode = buf.Next(1)[0]
	m.HardwareType = buf.Next(1)[0]
	m.HardwareLength = buf.Next(1)[0]
	m.Hops = buf.Next(1)[0]
	m.TransactionID = BytesToUint32(buf.Next(4))
	m.SecondsElapsed = BytesToUint16(buf.Next(2))
	m.Flags = BytesToUint16(buf.Next(2))
	m.ClientIP = buf.Next(4)
	m.YourIP = buf.Next(4)
	m.NextServerIP = buf.Next(4)
	m.RelayAgentIP = buf.Next(4)
	clientMac := ClientHardware{}
	clientMac.HardwareAddress = buf.Next(6)
	clientMac.HardwareAddressPadding = buf.Next(10)
	m.ClientMAC = clientMac
	m.ServerHostName = buf.Next(64)
	m.BootFile = buf.Next(128)
	m.MagicCookie = buf.Next(4)
	//decode options
	var options []OptionInter
	option := buf.Next(1)[0]
	for {
		switch option {
		case 1:
			options = append(options, Option1{}.Decode(buf.Next(5)))
		case 3:
			o := Option3{}
			o.Code = 3
			o.Length = buf.Next(1)[0]
			options = append(options, o.Decode(o.Length, buf.Next(int(o.Length))))
		case 6:
			o := Option6{}
			o.Code = 6
			o.Length = buf.Next(1)[0]
			options = append(options, o.Decode(o.Length, buf.Next(int(o.Length))))
		case 51:
			options = append(options, Option51{}.Decode(buf.Next(5)))
		case 53:
			option53 := Option53{}.Decode(buf.Next(2))
			options = append(options, option53)
			m.MessageType = option53.MessageType
		case 54:
			options = append(options, Option54{}.Decode(buf.Next(5)))
		case 58:
			options = append(options, Option58{}.Decode(buf.Next(5)))
		case 59:
			options = append(options, Option59{}.Decode(buf.Next(5)))
		case 61:
			length := buf.Next(1)[0]
			o := Option61{}.Decode(buf.Next(int(length)))
			o.Length = length
			options = append(options, o)
		case 108:
			options = append(options, Option108{}.Decode(buf.Next(5)))
		case 138:
			length := buf.Next(1)[0]
			o := Option138{}.Decode(buf.Next(int(length)))
			o.Length = length
			options = append(options, o)
		case 255:
			options = append(options, Option255{}.Decode([]byte{255}))
		default:
			m.Options = options
			return
		}
		option = buf.Next(1)[0]
	}
}

func (m *Message) String() string {
	var buf bytes.Buffer
	buf.WriteString("Message Type:")
	buf.WriteString(hex.EncodeToString(Uint8ToBytes(m.OpCode)))
	buf.WriteString("\n")
	buf.WriteString("Hardware Type:")
	buf.WriteString(hex.EncodeToString(Uint8ToBytes(m.HardwareType)))
	buf.WriteString("\n")
	buf.WriteString("Hardware address length:")
	buf.WriteString(hex.EncodeToString(Uint8ToBytes(m.HardwareLength)))
	buf.WriteString("\n")
	buf.WriteString("Hops:")
	buf.WriteString(hex.EncodeToString(Uint8ToBytes(m.Hops)))
	buf.WriteString("\n")
	buf.WriteString("Transaction ID:")
	buf.WriteString(hex.EncodeToString(Uint32ToBytes(m.TransactionID)))
	buf.WriteString("\n")
	buf.WriteString("Seconds elapsed:")
	buf.WriteString(hex.EncodeToString(Uint16ToBytes(m.SecondsElapsed)))
	buf.WriteString("\n")
	buf.WriteString("Bootp flags:")
	buf.WriteString(hex.EncodeToString(Uint16ToBytes(m.Flags)))
	buf.WriteString("\n")
	buf.WriteString("Client IP address:")
	buf.WriteString(net.IP(m.ClientIP).String())
	buf.WriteString("\n")
	buf.WriteString("Your (client) IP address:")
	buf.WriteString(net.IP(m.YourIP).String())
	buf.WriteString("\n")
	buf.WriteString("Next server IP address:")
	buf.WriteString(net.IP(m.NextServerIP).String())
	buf.WriteString("\n")
	buf.WriteString("Relay agent IP address:")
	buf.WriteString(net.IP(m.RelayAgentIP).String())
	buf.WriteString("\n")
	buf.WriteString("Client MAC address:")
	buf.WriteString(net.HardwareAddr(m.ClientMAC.HardwareAddress).String())
	buf.WriteString("\n")
	buf.WriteString("Client MAC address padding:")
	buf.WriteString(hex.EncodeToString(m.ClientMAC.HardwareAddressPadding))
	buf.WriteString("\n")
	buf.WriteString("Magic cookie:")
	buf.WriteString(hex.EncodeToString(m.MagicCookie))
	buf.WriteString("\n")
	for _, option := range m.Options {
		buf.WriteString(option.String())
		buf.WriteString("\n")
	}

	return buf.String()
}
