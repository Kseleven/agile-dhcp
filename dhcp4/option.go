package dhcp4

import (
	"bytes"
	"net"
	"strconv"
)

type OptionInter interface {
	Encode() []byte
	String() string
	GetCode() uint8
}

type MessageType uint8

const (
	MessageTypeDiscover MessageType = iota + 1
	MessageTypeOffer
	MessageTypeRequest
	MessageTypeDecline
	MessageTypeAck
	MessageTypeNak
	MessageTypeRelease
)

func (o MessageType) String() string {
	switch o {
	case MessageTypeDiscover:
		return "Discover"
	case MessageTypeOffer:
		return "Offer"
	case MessageTypeRequest:
		return "Request"
	case MessageTypeDecline:
		return "Decline"
	case MessageTypeAck:
		return "ACK"
	case MessageTypeNak:
		return "NAK"
	case MessageTypeRelease:
		return "Release"
	default:
		return ""
	}
}

// Option1 Subnet Mask
//The subnet mask option specifies the client's subnet mask as per RFC 950 [5].
//If both the subnet mask and the router option are specified in a DHCP reply, the subnet mask option MUST be first.
//The code for the subnet mask option is 1, and its length is 4 octets.
//    Code   Len        Subnet Mask
//   +-----+-----+-----+-----+-----+-----+
//   |  1  |  4  |  m1 |  m2 |  m3 |  m4 |
//   +-----+-----+-----+-----+-----+-----+
type Option1 struct {
	Code       uint8
	Length     uint8
	SubnetMask []byte
}

func (o Option1) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.SubnetMask...)
}

func (o Option1) Decode(b []byte) Option1 {
	o.Code = 1
	o.Length = b[0]
	o.SubnetMask = b[1:]
	return o
}

func (o Option1) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Subnet Mask:")
	buf.WriteString(net.IPMask(o.SubnetMask).String())

	return buf.String()
}

func (o Option1) GetCode() uint8 {
	return o.Code
}

//Option3 Router Option
//The router option specifies a list of IP addresses for routers on the
//   client's subnet.  Routers SHOULD be listed in order of preference.
//   The code for the router option is 3.  The minimum length for the
//   router option is 4 octets, and the length MUST always be a multiple of 4.
//    Code   Len         Address 1               Address 2
//   +-----+-----+-----+-----+-----+-----+-----+-----+--
//   |  3  |  n  |  a1 |  a2 |  a3 |  a4 |  a1 |  a2 |  ...
//   +-----+-----+-----+-----+-----+-----+-----+-----+--
type Option3 struct {
	Code    uint8
	Length  uint8
	Routers [][]byte
}

func (o Option3) Encode() []byte {
	b := []byte{o.Code, o.Length}
	for _, router := range o.Routers {
		b = append(b, router...)
	}
	return b
}

func (o Option3) Decode(length uint8, b []byte) Option3 {
	o.Code = 3
	for i := 0; i < int(length); i += 4 {
		o.Routers = append(o.Routers, b[i:i+4])
	}
	return o
}

func (o Option3) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Routers:")
	for _, router := range o.Routers {
		buf.WriteString(net.IP(router).String())
		buf.WriteString(" ")
	}

	return buf.String()
}

func (o Option3) GetCode() uint8 {
	return o.Code
}

//Option6 Domain Name Server Option
//The domain name server option specifies a list of Domain Name System
//   (STD 13, RFC 1035 [8]) name servers available to the client.  Servers
//   SHOULD be listed in order of preference.
//   The code for the domain name server option is 6.  The minimum length
//   for this option is 4 octets, and the length MUST always be a multiple of 4.
//
//    Code   Len         Address 1               Address 2
//   +-----+-----+-----+-----+-----+-----+-----+-----+--
//   |  6  |  n  |  a1 |  a2 |  a3 |  a4 |  a1 |  a2 |  ...
//   +-----+-----+-----+-----+-----+-----+-----+-----+--
type Option6 struct {
	Code              uint8
	Length            uint8
	DomainNameServers [][]byte
}

func (o Option6) Encode() []byte {
	b := []byte{o.Code, o.Length}
	for _, domainServer := range o.DomainNameServers {
		b = append(b, domainServer...)
	}
	return b
}

func (o Option6) Decode(length uint8, b []byte) Option6 {
	o.Code = 6
	for i := 0; i < int(length); i += 4 {
		o.DomainNameServers = append(o.DomainNameServers, b[i:i+4])
	}
	return o
}

func (o Option6) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Routers:")
	for _, domainServer := range o.DomainNameServers {
		buf.WriteString(net.IP(domainServer).String())
		buf.WriteString(" ")
	}

	return buf.String()
}

func (o Option6) GetCode() uint8 {
	return o.Code
}

//Option12 Host Name Option
//The code for this option is 12, and its minimum length is 1.
//
//    Code   Len                 Host Name
//   +-----+-----+-----+-----+-----+-----+-----+-----+--
//   |  12 |  n  |  h1 |  h2 |  h3 |  h4 |  h5 |  h6 |  ...
//   +-----+-----+-----+-----+-----+-----+-----+-----+--
type Option12 struct {
	Code     uint8
	Length   uint8
	HostName []byte
}

func (o Option12) GetCode() uint8 {
	return o.Code
}

func GenOption12(hostName string) Option12 {
	o := Option12{Code: 12}
	var buf bytes.Buffer
	buf.WriteString(hostName)
	o.HostName = buf.Bytes()
	o.Length = uint8(len(o.HostName))
	return o
}

func (o Option12) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.HostName...)
}

func (o Option12) Decode(b []byte) Option12 {
	var buf bytes.Buffer
	buf.Grow(len(b))
	buf.Write(b)

	o.Code = 12
	o.Length = buf.Next(1)[0]
	o.HostName = buf.Next(int(o.Length))
	return o
}

func (o Option12) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Host Name:")
	buf.Write(o.HostName)

	return buf.String()
}

//Option50 Requested IP Address
//This option is used in a client request (DHCPDISCOVER) to allow the
//   client to request that a particular IP address be assigned.
//   The code for this option is 50, and its length is 4.
//
//    Code   Len          Address
//   +-----+-----+-----+-----+-----+-----+
//   |  50 |  4  |  a1 |  a2 |  a3 |  a4 |
//   +-----+-----+-----+-----+-----+-----+
type Option50 struct {
	Code    uint8
	Length  uint8
	Address []byte
}

func (o Option50) GetCode() uint8 {
	return o.Code
}

func GenOption50(address []byte) Option50 {
	o := Option50{Code: 50}
	o.Length = uint8(4)
	o.Address = address
	return o
}

func (o Option50) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.Address...)
}

func (o Option50) Decode(b []byte) Option50 {
	var buf bytes.Buffer
	buf.Grow(len(b))
	buf.Write(b)

	o.Code = 50
	o.Length = buf.Next(1)[0]
	o.Address = buf.Next(int(o.Length))
	return o
}

func (o Option50) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Request IP Address:")
	buf.WriteString(net.IP(o.Address).String())

	return buf.String()
}

//Option51 IP Address Lease Time
//This option is used in a client request (DHCPDISCOVER or DHCPREQUEST)
//   to allow the client to request a lease time for the IP address.  In a
//   server reply (DHCPOFFER), a DHCP server uses this option to specify
//   the lease time it is willing to offer.
//
//   The time is in units of seconds, and is specified as a 32-bit
//   unsigned integer.
//
//   The code for this option is 51, and its length is 4.
//
//    Code   Len         Lease Time
//   +-----+-----+-----+-----+-----+-----+
//   |  51 |  4  |  t1 |  t2 |  t3 |  t4 |
//   +-----+-----+-----+-----+-----+-----+
type Option51 struct {
	Code      uint8
	Length    uint8
	LeaseTime []byte //32-bit
}

func GenOption51(t uint32) Option51 {
	return Option51{Code: 51, Length: 4, LeaseTime: Uint32ToBytes(t)}
}

func (o Option51) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.LeaseTime...)
}

func (o Option51) Decode(b []byte) Option51 {
	o.Code = 51
	o.Length = b[0]
	o.LeaseTime = b[1:]
	return o
}

func (o Option51) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" IP Address Lease Time:")
	buf.WriteString(strconv.FormatUint(uint64(BytesToUint32(o.LeaseTime)), 10))

	return buf.String()
}

func (o Option51) GetCode() uint8 {
	return o.Code
}

// Option53 DHCP Message Type(3 octets)
//Value   Message Type
//-----   ------------
//1     DHCPDISCOVER
//2     DHCPOFFER
//3     DHCPREQUEST
//4     DHCPDECLINE
//5     DHCPACK
//6     DHCPNAK
//7     DHCPRELEASE
//Code   Len  Type
//+-----+-----+-----+
//|  53 |  1  | 1-7 |
//+-----+-----+-----+
type Option53 struct {
	Code        uint8
	Length      uint8
	MessageType MessageType
}

func (o Option53) GetCode() uint8 {
	return o.Code
}

func GenOption53(c MessageType) Option53 {
	return Option53{Code: 53, Length: 1, MessageType: c}
}

func (o Option53) Encode() []byte {
	return []byte{o.Code, o.Length, uint8(o.MessageType)}
}

func (o Option53) Decode(b []byte) Option53 {
	o.Code = 53
	o.Length = b[0]
	o.MessageType = MessageType(b[1])
	return o
}

func (o Option53) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" DHCP:")
	buf.WriteString(o.MessageType.String())
	return buf.String()
}

// Option54 Server Identifier
//The code for this option is 54, and its length is 4.
//
//    Code   Len            Address
//   +-----+-----+-----+-----+-----+-----+
//   |  54 |  4  |  a1 |  a2 |  a3 |  a4 |
//   +-----+-----+-----+-----+-----+-----+
type Option54 struct {
	Code             uint8
	Length           uint8
	ServerIdentifier []byte
}

func (o Option54) GetCode() uint8 {
	return o.Code
}

func (o Option54) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.ServerIdentifier...)
}

func (o Option54) Decode(b []byte) Option54 {
	var buf bytes.Buffer
	buf.Grow(len(b))
	buf.Write(b)

	o.Code = 54
	o.Length = buf.Next(1)[0]
	o.ServerIdentifier = buf.Next(4)
	return o
}

func (o Option54) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Server Identifier:")
	buf.WriteString(net.IP(o.ServerIdentifier).String())

	return buf.String()
}

// Option55 Parameter Request List
//Code   Len   Option Codes
//   +-----+-----+-----+-----+---
//   |  55 |  n  |  c1 |  c2 | ...
//   +-----+-----+-----+-----+---
//The code for this option is 55.  Its minimum length is 1.
type Option55 struct {
	Code       uint8
	Length     uint8
	Parameters []byte
}

func (o Option55) GetCode() uint8 {
	return o.Code
}

func GenOption55() Option55 {
	//option1: Subnet Mask
	//option3: Router
	//option6: Domain Name Server
	//option15: Domain Name
	//option46: NetBIOS over TCP/IP Node Type
	//option95: LDAP
	//option108: IPv6-Only Preferred
	//option114: DHCP Captive-Portal(URL)
	//option118: Subnet Selection Option
	//option119: DNS Domain Search List
	//option121: Classless Static Route
	//option252: Private/Proxy autodiscovery
	var parameters = []byte{1, 3, 6, 15, 46, 108, 114, 119, 121, 252}
	return Option55{Code: 55, Length: uint8(len(parameters)), Parameters: parameters}
}

func (o Option55) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteByte(o.Code)
	buf.WriteByte(o.Length)
	buf.Write(o.Parameters)
	return buf.Bytes()
}

func (o Option55) Decode(b []byte) Option55 {
	var buf bytes.Buffer
	buf.Grow(len(b))
	buf.Write(b)

	o.Code = 55
	o.Length = buf.Next(1)[0]
	o.Parameters = append(o.Parameters, buf.Next(buf.Len())...)
	return o
}

func (o Option55) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Parameter Request List Item:")
	for _, parameter := range o.Parameters {
		buf.WriteString(strconv.FormatUint(uint64(parameter), 10))
		buf.WriteString(" ")
	}
	return buf.String()
}

//Option57 Maximum DHCP Message Size
//This option specifies the maximum length DHCP message that it is
//   willing to accept.  The length is specified as an unsigned 16-bit
//   integer.  A client may use the maximum DHCP message size option in
//   DHCPDISCOVER or DHCPREQUEST messages, but should not use the option in DHCPDECLINE messages.
//The code for this option is 57, and its length is 2.  The minimum
//   legal value is 576 octets.
//
//    Code   Len     Length
//   +-----+-----+-----+-----+
//   |  57 |  2  |  l1 |  l2 |
//   +-----+-----+-----+-----+
type Option57 struct {
	Code               uint8
	Length             uint8
	MaximumMessageSize []byte
}

func GenOption57(t uint16) Option57 {
	return Option57{Code: 57, Length: 2, MaximumMessageSize: Uint16ToBytes(t)}
}

func (o Option57) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.MaximumMessageSize...)
}

func (o Option57) Decode(b []byte) Option57 {
	o.Length = b[0]
	o.MaximumMessageSize = b[1:]
	return o
}

func (o Option57) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Maximum DHCP Message Size:")
	buf.WriteString(strconv.FormatUint(uint64(BytesToUint16(o.MaximumMessageSize)), 10))

	return buf.String()
}

func (o Option57) GetCode() uint8 {
	return o.Code
}

//Option58 Renewal (T1) Time Value
//This option specifies the time interval from address assignment until
//   the client transitions to the RENEWING state.
//   The value is in units of seconds, and is specified as a 32-bit
//   unsigned integer.
//
//   The code for this option is 58, and its length is 4.
//
//    Code   Len         T1 Interval
//   +-----+-----+-----+-----+-----+-----+
//   |  58 |  4  |  t1 |  t2 |  t3 |  t4 |
//   +-----+-----+-----+-----+-----+-----+
type Option58 struct {
	Code        uint8
	Length      uint8
	RenewalTime []byte //32-bit
}

func (o Option58) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.RenewalTime...)
}

func (o Option58) Decode(b []byte) Option58 {
	o.Code = 58
	o.Length = b[0]
	o.RenewalTime = b[1:]
	return o
}

func (o Option58) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Renewal Time Value:")
	buf.WriteString(strconv.FormatUint(uint64(BytesToUint32(o.RenewalTime)), 10))

	return buf.String()
}

func (o Option58) GetCode() uint8 {
	return o.Code
}

//Option59 Rebinding (T2) Time Value
//This option specifies the time interval from address assignment until
//   the client transitions to the REBINDING state.
//
//   The value is in units of seconds, and is specified as a 32-bit
//   unsigned integer.
//
//   The code for this option is 59, and its length is 4.
//
//    Code   Len         T2 Interval
//   +-----+-----+-----+-----+-----+-----+
//   |  59 |  4  |  t1 |  t2 |  t3 |  t4 |
//   +-----+-----+-----+-----+-----+-----+
type Option59 struct {
	Code          uint8
	Length        uint8
	RebindingTime []byte //32-bit
}

func (o Option59) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.RebindingTime...)
}

func (o Option59) Decode(b []byte) Option59 {
	o.Code = 59
	o.Length = b[0]
	o.RebindingTime = b[1:]
	return o
}

func (o Option59) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Rebinding Time Value:")
	buf.WriteString(strconv.FormatUint(uint64(BytesToUint32(o.RebindingTime)), 10))

	return buf.String()
}

func (o Option59) GetCode() uint8 {
	return o.Code
}

//Option61 Client-identifier
//This option is used by DHCP clients to specify their unique
//   identifier.  DHCP servers use this value to index their database of
//   address bindings.  This value is expected to be unique for all
//   clients in an administrative domain.
//Identifiers consist of a type-value pair
//   It is expected that this field will typically contain a hardware type
//   and hardware address, but this is not required.  Current legal values
//   for hardware types are defined in [22].
//The code for this option is 61, and its minimum length is 2.
//
//   Code   Len   Type  Client-Identifier
//   +-----+-----+-----+-----+-----+---
//   |  61 |  n  |  t1 |  i1 |  i2 | ...
//   +-----+-----+-----+-----+-----+---
type Option61 struct {
	Code             uint8
	Length           uint8
	HardwareType     uint8 //Hardware type Ethernet
	ClientIdentifier []byte
}

func GenOption61(mac []byte) Option61 {
	return Option61{Code: 61, Length: 7, HardwareType: 1, ClientIdentifier: mac}
}

func (o Option61) Encode() []byte {
	return append([]byte{o.Code, o.Length, o.HardwareType}, o.ClientIdentifier...)
}

func (o Option61) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" Hardware type Ethernet:")
	buf.WriteString(strconv.FormatUint(uint64(o.HardwareType), 10))
	buf.WriteString(" Client MAC address:")
	buf.WriteString(net.HardwareAddr(o.ClientIdentifier).String())

	return buf.String()
}

func (o Option61) GetCode() uint8 {
	return o.Code
}

func (o Option61) Decode(b []byte) Option61 {
	o.Code = 61
	o.HardwareType = b[0]
	o.ClientIdentifier = b[1:]
	return o
}

//Option108 IPv6-Only Preferred Option
//Code:
//8-bit identifier of the IPv6-Only Preferred option code as assigned by IANA: 108.
//The client includes the Code in the Parameter Request List in DHCPDISCOVER and DHCPREQUEST messages as described in Section 3.2.
//Length:
//8-bit unsigned integer. The length of the option, excluding the Code and Length Fields.
//The server MUST set the length field to 4. The client MUST ignore the IPv6-Only Preferred option if the length field value is not 4.
//Value:
//32-bit unsigned integer. The number of seconds for which the client should disable DHCPv4 (V6ONLY_WAIT configuration variable).
//If the server pool is explicitly configured with a V6ONLY_WAIT timer,
//the server MUST set the field to that configured value. Otherwise,
//the server MUST set it to zero. The client MUST process that field as described in Section 3.2.
//
//The client never sets this field, as it never sends the full option
//but includes the option code in the Parameter Request List as described in Section 3.2.
// 0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Code      |   Length      |           Value               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Value (cont.)         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type Option108 struct {
	Code   uint8
	Length uint8
	Value  []byte
}

func (o Option108) Encode() []byte {
	return append([]byte{o.Code, o.Length}, o.Value...)
}

func (o Option108) Decode(b []byte) Option108 {
	o.Code = 108
	o.Length = b[0]
	o.Value = b[1:]
	return o
}

func (o Option108) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Length:")
	buf.WriteByte(o.Length)
	buf.WriteString(strconv.FormatUint(uint64(o.Length), 10))
	buf.WriteString(" IPv6-Only Preferred:")
	buf.WriteString(strconv.FormatUint(uint64(BytesToUint32(o.Value)), 10))

	return buf.String()
}

func (o Option108) GetCode() uint8 {
	return o.Code
}

// Option255 End Option
//The end option marks the end of valid information in the vendor
//   field.  Subsequent octets should be filled with pad options.
//   The code for the end option is 255, and its length is 1 octet.
type Option255 struct {
	Code uint8
}

func (o Option255) GetCode() uint8 {
	return o.Code
}

func GenOption255() Option255 {
	return Option255{Code: 0xff}
}

func (o Option255) Encode() []byte {
	return []byte{o.Code}
}

func (o Option255) Decode(b []byte) Option255 {
	o.Code = b[0]
	return o
}

func (o Option255) String() string {
	var buf bytes.Buffer
	buf.WriteString("Option:(")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	buf.WriteString(")")
	buf.WriteString(" Option End:")
	buf.WriteString(strconv.FormatUint(uint64(o.Code), 10))
	return buf.String()
}
