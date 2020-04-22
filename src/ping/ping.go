package ping

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gvso/internship-application-systems/src/icmp"
)

var ip4addr = "0.0.0.0"
var ip6addr = "::"

// Ping sends a ICMP echo request.
func Ping(addr string) (*net.IPAddr, time.Duration, error) {
	c, err := icmp.ListenPacket(icmp.IPv4Type, ip4addr)

	if err != nil {
		return nil, 0, err
	}

	defer c.Close()

	// Resolve any DNS (if used) and get the real IP of the destination.
	dst, err := net.ResolveIPAddr(c.Protocol, addr)
	if err != nil {
		return nil, 0, err
	}

	echo := icmp.Body{
		ID:   1 & 0xffff,
		Seq:  9 & 0xffff,
		Data: []byte("TEST"),
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
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(body))
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	n, _, err = c.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}

	// Pack it up boys, we're done here
	msg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return dst, 0, err
	}

	log.Println(msg)

	if msg.Type == 0 {
		duration := time.Since(start)
		log.Println(duration)
	}

	return dst, 0, nil
}
