package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"net/http"
	"os"
	"time"
)

var ping_ip = "8.8.8.8"
var http_url = "http://www.google.com"

var bad_ping = 300 * time.Millisecond

func main() {
	ping_rtt := time.Duration(0)
	http_ok := false

	http_transport := http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
	}

	client := http.Client{
		Transport: &http_transport,
	}

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ping_ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ping_rtt = rtt
	}
	p.OnIdle = func() {
	}
	p.RunLoop()

	time.Sleep(time.Second * 2)
	p.Done()
	if ping_rtt == 0 {
		_, err = client.Head(http_url)
		if err == nil {
			http_ok = true
		}
		if http_ok {
			print_limit()
		} else {
			print_dead()
		}
	} else {
		print_ping(ping_rtt)
		// all fine, sleep for a bit
	}
	os.Exit(0)
}

func print_ping(ping time.Duration) {
	fmt.Printf("RTT: %v\n", ping)
	fmt.Printf("RTT: %v\n", ping)
	if ping > bad_ping {
		fmt.Println("#FF9900")
	} else {
		fmt.Println("#00FF00")
	}
}

func print_limit() {
	fmt.Printf("HTTP!\n")
	fmt.Printf("HTTP!\n")
	fmt.Println("#FFCC00")
}

func print_dead() {
	fmt.Printf("HTTP!\n")
	fmt.Printf("HTTP!\n")
	fmt.Println("#FF0000")
}
