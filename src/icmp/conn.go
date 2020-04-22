package icmp

import (
	"errors"
	"net"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var errInvalid = errors.New("Invalid connection")

// PacketConn defines a ICMP connection.
type PacketConn struct {
	Protocol string
	c        net.PacketConn
	v4       *ipv4.PacketConn
	v6       *ipv6.PacketConn
}

func (c *PacketConn) ok() bool {
	return c != nil && c.c != nil
}

// Close closes the endpoint.
func (c *PacketConn) Close() error {
	if !c.ok() {
		return errInvalid
	}

	return c.c.Close()
}

// WriteTo writes the ICMP message b to dst.
func (c *PacketConn) WriteTo(b []byte, dst net.Addr) (int, error) {
	if !c.ok() {
		return 0, errInvalid
	}

	return c.c.WriteTo(b, dst)
}

// ReadFrom reads an ICMP message.
func (c *PacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	if !c.ok() {
		return 0, nil, errInvalid
	}

	return c.c.ReadFrom(b)
}

// SetReadDeadline sets the read deadline associated with the connection.
func (c *PacketConn) SetReadDeadline(t time.Time) error {
	if !c.ok() {
		return errInvalid
	}

	return c.c.SetReadDeadline(t)
}
