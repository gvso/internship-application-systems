package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gvso/internship-application-systems/src/icmp"
)

// BuffSize is the size of the receiving buffer.
const BuffSize = 1024

const ip4addr = "0.0.0.0"
const ip6addr = "::"

// makePing sends and receives the ICMP messages.
func makePing(c *icmp.PacketConn, dst *net.IPAddr, seq int) (*icmp.Message, time.Duration, error) {
	echo := icmp.Body{
		ID:   os.Getpid() & 0xffff,
		Seq:  seq & 0xffff,
		Data: []byte(""),
	}

	request := icmp.NewEchoRequest(c.Protocol, echo)

	body, err := request.Marshal()

	if err != nil {
		panic(err)
	}

	start := time.Now()
	n, err := c.WriteTo(body, dst)
	if err != nil {
		panic(err)
	} else if n != len(body) {
		return nil, 0, fmt.Errorf("got %v; want %v", n, len(body))
	}

	// Sets a deadline to wait for a reply.
	err = c.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return nil, 0, err
	}

	// Wait for a reply
	reply := make([]byte, BuffSize)
	n, _, err = c.ReadFrom(reply)
	if err != nil {
		return nil, 0, err
	}

	// Parse the ICMP message.
	msg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return msg, 0, err
	}

	duration := time.Since(start)

	return msg, duration, nil
}

// Ping sends a ICMP echo request.
func Ping(addr string, seq int) (*Response, error) {
	c, err := icmp.ListenPacket(icmp.IPv4Type, ip4addr)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	var resp Response

	// Resolve any DNS (if used) and get the IP and CNAME of the destination.
	cname, _ := net.LookupCNAME(addr)
	if err != nil {
		return nil, err
	}
	resp.CName = cname

	dst, err := net.ResolveIPAddr(c.Protocol, addr)
	if err != nil {
		return nil, err
	}
	resp.Dest = dst

	msg, duration, err := makePing(c, dst, seq)
	if err != nil {
		return nil, err
	}

	resp.Duration = duration.Round(time.Microsecond)
	resp.Type = msg.Type
	resp.Seq = msg.Body.Seq

	switch msg.Type {
	case icmp.ICMPTypeEchoReply:
		return &resp, nil
	default:
		resp.Type = -1
		return &resp, fmt.Errorf("got %+v from %v", msg, addr)
	}
}
