package main

import (
	"fmt"
	"github.com/danieldin95/lightstar/pkg/libstar"
	"io"
	"net"
)

func copyBytes(dir string, dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, 1024*32)
	for {
		nr, er := src.Read(buf)
		libstar.Info("Server.HandleWebsockify %s % x", dir, buf[:nr])
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:5900")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		target, err := net.Dial("tcp", "192.168.4.249:5900")
		if err != nil {
			libstar.Error("main.Dial %s", err)
		}
		//target, err := websocket.Dial("ws://192.168.209.141:10080/websockify", "", "http://192.168.209.141:6080/")

		go copyBytes("from target", conn, target)
		go copyBytes("from server", target, conn)
	}
}
