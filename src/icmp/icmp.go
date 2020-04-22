package icmp

import (
	"net"
	"strconv"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Internet Assigned Numbers Authority (IANA) protocol numbers.
const (
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
)

// Protocol type as strings.
const (
	IPv4Type = "ip4"
	IPv6Type = "ip6"
)

// ICMP Types
const (
	ICMPTypeEchoReply   = 0
	ICMPTypeEchoRequest = 8
)

// ListenPacket starts listening for ICMP messages.
func ListenPacket(proto, address string) (*PacketConn, error) {

	var c net.PacketConn
	network := proto + ":"

	switch proto {
	case IPv6Type:
		network += strconv.Itoa(ProtocolIPv6ICMP)
	default:
		network += strconv.Itoa(ProtocolICMP)
	}

	c, err := net.ListenPacket(network, address)

	if err != nil {
		return nil, err
	}

	switch proto {
	case IPv6Type:
		return &PacketConn{Protocol: "ip6", c: c, v6: ipv6.NewPacketConn(c)}, nil
	default:
		return &PacketConn{Protocol: "ip4", c: c, v4: ipv4.NewPacketConn(c)}, nil
	}
}
