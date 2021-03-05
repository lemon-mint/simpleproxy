package proxylib

import (
	"fmt"
	"net"
	"time"
)

//Proxy : proxy object
type Proxy struct {
	Protocol    string
	ListenAddr  string
	Destination string

	Unit            int
	UseDelay        bool
	Delay           time.Duration
	DebugPrint      bool
	ConnectionPrint bool
}

func (p *Proxy) handleconn(sconn net.Conn) {
	dconn, err := net.Dial(p.Protocol, p.Destination)
	if err != nil {
		sconn.Close()
		return
	}
	go p.pipe(sconn, dconn, "src -> dst")
	go p.pipe(dconn, sconn, "dst -> src")
}

func (p *Proxy) pipe(src, dst net.Conn, id string) {
	buf := make([]byte, p.Unit)
	defer src.Close()
	defer dst.Close()
	for {
		n, err := src.Read(buf[:p.Unit])
		if err != nil {
			return
		}
		if p.UseDelay {
			time.Sleep(p.Delay)
		}
		_, err = dst.Write(buf[:n])
		if err != nil {
			return
		}
		if p.DebugPrint {
			fmt.Println(id, "Proxyed", n, "bytes.")
		}
	}
}

//Serve : start proxy server
func (p *Proxy) Serve() error {
	l, err := net.Listen(p.Protocol, p.ListenAddr)
	if err != nil {
		return err
	}
	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			break
		}
		if p.ConnectionPrint {
			fmt.Println("connected", conn.LocalAddr(), conn.RemoteAddr())
		}
		go p.handleconn(conn)
	}
	return err
}
