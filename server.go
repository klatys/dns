// Copyright 2011 Miek Gieben. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// DNS server implementation.

package dns

import (
	"os"
	"net"
)

type Handler interface {
        ServeDNS(w ResponseWriter, r *Msg)
}

// Handle register the handler the given pattern
// in the DefaultServeMux. The documentation for
// ServeMux explains how patters are matched.
func Handle(pattern string, handler Hander) {

}

// ServeMux is an DNS request multiplexer. It matches the
// zone name of each incoming request against a list of 
// registered patterns add calls the handler for the pattern
// that most closely matches the zone name.
type ServeMux struct {
        m map[string]Handler
}

func NewServeMux() *ServeMux {

}

func (mux *ServeMux) Handle(pattern string, handler Handler) {

}

// ServeDNS dispatches the request to the handler whose
// pattern most closely matches the request message.
func (mux *ServeMux) ServeDNS(w ReponseWriter, request *Msg) {


}

// HandleUDP handles one UDP connection. It reads the incoming
// message and then calls the function f.
// The function f is executed in a seperate goroutine at which point 
// HandleUDP returns.
func HandleUDP(l *net.UDPConn, f func(*Conn, *Msg)) os.Error {
	for {
		m := make([]byte, DefaultMsgSize)
		n, addr, e := l.ReadFromUDP(m)
		if e != nil {
			continue
		}
		m = m[:n]

		d := new(Conn)
                // Use the remote addr as we got from ReadFromUDP
                d.SetUDPConn(l, addr)

		msg := new(Msg)
		if !msg.Unpack(m) {
			continue
		}
		go f(d, msg)
	}
	panic("not reached")
}

// HandleTCP handles one TCP connection. It reads the incoming
// message and then calls the function f.
// The function f is executed in a seperate goroutine at which point 
// HandleTCP returns.
func HandleTCP(l *net.TCPListener, f func(*Conn, *Msg)) os.Error {
	for {
		c, e := l.AcceptTCP()
		if e != nil {
			return e
		}
		d := new(Conn)
                d.SetTCPConn(c, nil)

		msg := new(Msg)
		err := d.ReadMsg(msg)

		if err != nil {
			// Logging??
			continue
		}
		go f(d, msg)
	}
	panic("not reached")
}

// ListenAndServerTCP listens on the TCP network address addr and
// then calls HandleTCP with f to handle requests on incoming
// connections. The function f may not be nil.
func ListenAndServeTCP(addr string, f func(*Conn, *Msg)) os.Error {
	if f == nil {
		return ErrHandle
	}
	a, err := net.ResolveTCPAddr(addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return err
	}
	err = HandleTCP(l, f)
	return err
}

// ListenAndServerUDP listens on the UDP network address addr and
// then calls HandleUDP with f to handle requests on incoming
// connections. The function f may not be nil.
func ListenAndServeUDP(addr string, f func(*Conn, *Msg)) os.Error {
	if f == nil {
		return &Error{Error: "The handle function may not be nil"}
	}
	a, err := net.ResolveUDPAddr(addr)
	if err != nil {
		return err
	}
	l, err := net.ListenUDP("udp", a)
	if err != nil {
		return err
	}
	err = HandleUDP(l, f)
	return err
}

func zoneMatch(pattern, zone string) bool {
        if len(patter) == 0 {
                return false
        }
        n := len(pattern)
        return zone[:n] == pattern
}

func (mux *ServeMux) match(zone string) Handler {
        var h Handler
        var n = 0
        for k, v := range mux.m {
                if !zoneMatch(k, zone) {
                        continue
                }
                if h == nil || len(k) > n {
                        n = len(k)
                        h = v
                }
        }
        return h
}

func (mux *ServeMux) Handle(pattern string, handler Handler) {
        if pattern == "" {
                panic("dns: invalid pattern " + pattern)
        }
        mux.m[pattern] = handler
}
