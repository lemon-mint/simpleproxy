package proxylib

import (
	"net"
	"time"
)

//Proxy : proxy object
type Proxy struct {
	Protocol    string
	ListenAddr  string
	Destination string

	Unit     int
	UseDelay bool
	Delay    time.Duration
}

func (p *Proxy) handleconn(sconn net.Conn) {
	dconn, err := net.Dial(p.Protocol, p.Destination)
	if err != nil {
		sconn.Close()
		dconn.Close()
		return
	}
	go p.pipe(sconn, dconn)
	go p.pipe(dconn, sconn)
}

func (p *Proxy) pipe(src, dst net.Conn) {
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
		go p.handleconn(conn)
	}
	return err
}
