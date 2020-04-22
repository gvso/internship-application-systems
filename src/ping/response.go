package ping

import (
	"net"
	"time"
)

// Response is the format of the ping response.
type Response struct {
	Dest     *net.IPAddr
	CName    string
	Duration time.Duration
	Type     int
	Seq      int
}
