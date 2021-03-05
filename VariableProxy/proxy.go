package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lemon-mint/simpleproxy/proxylib"
)

const bufsize = 32768

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 4 {
		return
	}
	proto := os.Args[1]
	fmt.Println("Protocol :", proto)
	dest := os.Args[2]
	fmt.Println("Destination :", dest)
	listenAddr := os.Args[3]
	fmt.Println("ListenAddr :", listenAddr)
	var UseDelay bool = false
	var Delay time.Duration
	if len(os.Args) == 4 {
		UseDelay = false
	} else {
		UseDelay = true
		sdelay := os.Args[4]
		dtime, err := strconv.Atoi(sdelay)
		if err != nil {
			log.Fatalln(err)
		}
		Delay = time.Millisecond * time.Duration(dtime)
		fmt.Println("Delay : true")
	}
	p := &proxylib.Proxy{
		Protocol:        proto,
		ListenAddr:      listenAddr,
		Destination:     dest,
		UseDelay:        UseDelay,
		Delay:           Delay,
		Unit:            bufsize,
		DebugPrint:      false,
		ConnectionPrint: true,
	}
	go func() {
		log.Fatalln(p.Serve())
	}()
	var latency int
	for {
		fmt.Print("latency (ms) >> ")
		fmt.Scanln(&latency)
		if latency == 0 {
			p.UseDelay = false
		} else {
			p.UseDelay = true
		}
		p.Delay = time.Millisecond * time.Duration(latency)
	}
}
