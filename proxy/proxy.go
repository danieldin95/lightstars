package proxy

import (
	"fmt"
	"github.com/danieldin95/lightstar/libstar"
	"golang.org/x/net/websocket"
	"io"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	START = 9000
	SIZE  = 1000
)

var Start = START

type Target struct {
	Name   string `json:"name"`
	Target string `json:"target"`
}

type Local struct {
	Tgt      Target
	Listen   string
	Listener net.Listener
	Client   *WsClient
}

func (l *Local) Initialize() *Local {
	for i := Start; i < SIZE+START; i++ {
		// random or sequence
		l.Listen = fmt.Sprintf("localhost:%d", i)
		listen, err := net.Listen("tcp", l.Listen)
		if err == nil {
			Start = i
			l.Listener = listen
			break
		}
	}
	if l.Listener != nil {
		libstar.Info("Local.Initialize %-15s %-20s on %-15s", l.Tgt.Name, l.Tgt.Target, l.Listen)
	}
	return l
}

func (l *Local) Start() {
	libstar.Debug("Local.Accept %s", l.Listen)
	if l.Listener == nil || l.Client == nil {
		libstar.Error("Local.Accept: Invalid Local")
	}

	defer l.Listener.Close()
	for {
		conn, err := l.Listener.Accept()
		if err != nil {
			libstar.Error("Local.Accept: %s", err)
			return
		}
		libstar.Info("Local.Accept %s", conn.RemoteAddr())
		go l.OnClient(conn)
	}
}

func (l *Local) OnClient(from net.Conn) {
	// proxy by websocket.
	ws := &WsClient{
		Auth: l.Client.Auth,
		Url:  l.Client.Url,
	}
	ws.Url += l.Tgt.Target
	ws.Initialize()
	to, err := ws.Dial()
	if err != nil {
		libstar.Error("Local.Socket dial %s", err)
		return
	}

	// wait exit.
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(from, to); err != nil {
			libstar.Error("Local.Handle copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(to, from); err != nil {
			libstar.Error("Local.Handle copy from local %s", err)
		}
	}()
	wait.Wait()
}

func (l *Local) Stop() {
	if l.Listener != nil {
		l.Listener.Close()
	}
}

type Proxy struct {
	Target []Target
	Listen map[string]*Local
	Client *WsClient
	Conn   *websocket.Conn
}

func (p *Proxy) Initialize() *Proxy {
	if p.Listen == nil {
		p.Listen = make(map[string]*Local, 32)
	}
	for _, tgt := range p.Target {
		if tgt.Target == "" || !strings.Contains(tgt.Target, ":") {
			continue
		}
		local := &Local{
			Tgt:    tgt,
			Client: p.Client,
		}
		local.Initialize()
		p.Listen[tgt.Target] = local
	}

	return p
}

func (p *Proxy) Start() {
	for _, local := range p.Listen {
		go local.Start()
	}
}

func (p *Proxy) Stop() {
	if p.Conn != nil {
		p.Conn.Close()
	}
	for tgt, local := range p.Listen {
		libstar.Info("Proxy.Stop %s", tgt)
		local.Stop()
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
