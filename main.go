package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gvso/internship-application-systems/src/ping"
)

var received = 0
var lost = 0

func printStatistics() {
	total := received + lost

	fmt.Printf("%d packets transmitted, %d received, %.0f%% packet loss\n", total, received, float32(lost)/float32(total)*100)
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Run the program using ./<binary> <host>")
	}

	// Validates the address.
	addr := os.Args[1]
	_, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		fmt.Printf("Invalid address %s", addr)
		os.Exit(1)
	}

	// Listens for OS signals.
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		printStatistics()
		os.Exit(1)
	}()

	fmt.Printf("PING %s\n", addr)

	for i := 0; true; i++ {
		res, err := ping.Ping(addr, i)

		if err != nil {
			if strings.Contains(err.Error(), "socket: operation not permitted") {
				log.Fatal("You need to run the program as sudo")
				os.Exit(1)
			}

			lost++
		} else {
			fmt.Printf("Message from %s (%v): icmp_seq=%d time=%v\n", res.CName, res.Dest, res.Seq, res.Duration)
			received++
		}

		time.Sleep(2 * time.Second)
	}

}
