package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"
)

var ping_ip = "8.8.8.8"
var http_url = "http://www.google.com/"

var bad_ping = 300 * time.Millisecond

func main() {
	ping_rtt := time.Duration(0)
	http_ok := false

	client := http.Client{
		Timeout: 7 * time.Second,
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
		time_start := time.Now()
		_, err = client.Head(http_url)
		http_rtt := time.Since(time_start)
		if err == nil {
			http_ok = true
		}
		if http_ok {
			print_http(http_rtt)
			// dont bother that url too much, if ping doesnt work but http doee
			// it usually means firewall eats it
			time.Sleep(time.Second * 120)
		} else {
			print_dead()
		}
	} else {
		print_ping(ping_rtt)
		// all fine, sleep for a bit
	}
	os.Exit(0)
}

func print_ping(rtt time.Duration) {
	fmt.Printf("RTT: %s\n", format_rtt(rtt))
	fmt.Printf("RTT: %s\n", format_rtt(rtt))
	if rtt > (bad_ping * 2) {
		fmt.Println("#FFCC00")
	} else if rtt > (bad_ping) {
		fmt.Println("#FFFF00")
	} else if rtt > (bad_ping / 2) {
		fmt.Println("#66FF00")
	} else if rtt > (bad_ping / 2) {
		fmt.Println("#33FF00")
	} else {
		fmt.Println("#00FF00")
	}
}

func print_http(rtt time.Duration) {
	fmt.Printf("HTTP: %s\n", format_rtt(rtt))
	fmt.Printf("HTTP: %s\n", format_rtt(rtt))
	bad_ping := bad_ping * 5 // more rtts to do http ping, especially if proxy is involved
	if rtt > (bad_ping * 2) {
		fmt.Println("#FF2200")
	} else if rtt > (bad_ping) {
		fmt.Println("#FFFF00")
	} else if rtt > (bad_ping / 2) {
		fmt.Println("#FFAAAA")
	} else if rtt > (bad_ping / 2) {
		fmt.Println("#FFAAFF")
	} else {
		fmt.Println("#AAAAFF")
	}
}

func print_dead() {
	fmt.Printf("DEAD!\n")
	fmt.Printf("DEAD!\n")
	fmt.Println("#FF0000")
}

func format_rtt(rtt time.Duration) string {
	str := fmt.Sprintf("%v", rtt)
	re := regexp.MustCompile(`\.(\d{3})\d*`)
	return re.ReplaceAllString(str, ".$1")
}
