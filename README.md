# Simple DHCP client

### Feature
* dhcp client4
  * option 1 (Subnet Mask)
  * option 3 (Router)
  * option 6 (Domain Name Server)
  * option 12 (Host Name)
  * option 51 (IP Address Lease Time)
  * option 53 (DHCP Message Type)
  * option 54 (Server Identifier)
  * option 55 (Parameter Request List)
  * option 57 (Maximum DHCP Message Size)
  * option 58 (Renewal (T1) Time Value)
  * option 59 (Rebinding (T2) Time Value)
  * option 61 (Client-identifier)
  * option 108 (IPv6-Only Preferred)
  * option 255 (End Option)
* dhcp client6 (going on)

### Usage
* run with source
```shell
git clone github.com/Kseleven/agile-dhcp
go run example/dhcp4/dhcp4.go -h test -m 00:00:00:00:00:01
```

* run with binary
```shell
git clone github.com/Kseleven/agile-dhcp
make 
./dhcp_client4  -h test -m 00:00:00:00:00:01
```

### Good luck
