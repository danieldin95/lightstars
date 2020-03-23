package http

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/service"
	"github.com/danieldin95/lightstar/storage"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/http"
	"sync"
)

type Download struct {
}

func (down Download) Router(router *mux.Router) {
	dir := http.Dir(storage.PATH.Root())
	files := http.StripPrefix("/ext/files/", http.FileServer(dir))
	router.PathPrefix("/ext/files/").Handler(files)
}

type WebSocket struct {
}

func (w WebSocket) Router(router *mux.Router) {
	router.Handle("/ext/websocket", websocket.Handler(w.Handle))
}

func (w WebSocket) GetRemote(id, name, typ string) string {
	libstar.Debug("WebSocket.GetRemote %s://%s@%s", typ, id, name)
	node := service.SERVICE.Zone.Get(name)
	if node == nil {
		libstar.Error("WebSocket.GetRemote %s", name)
		return ""
	}
	host := node.Hostname
	client := libstar.HttpClient{
		Url: node.Url + "/api/instance/" + id + "?format=schema",
		Auth: libstar.Auth{
			Type:     "basic",
			Password: node.Password,
			Username: node.Username,
		},
	}
	resp, err := client.Do()
	if err != nil {
		libstar.Error("WebSocket.GetRemote %s", err)
		return ""
	}
	inst := schema.Instance{}
	if err := libstar.GetJSON(resp.Body, &inst); err != nil {
		libstar.Error("WebSocket.GetRemote %s", name)
		return ""
	}
	port := inst.Vnc.Port
	if typ == "spice" {
		port = inst.Spice.Port
	}
	return host + ":" + port
}

func (w WebSocket) GetLocal(id, typ string) string {
	libstar.Debug("WebSocket.GetLocal %s://%s", typ, id)
	hyper, err := libvirtc.GetHyper()
	if err != nil {
		libstar.Error("WebSocket.GetLocal %s", err)
		return ""
	}
	dom, err := hyper.LookupDomainByUUIDString(id)
	if err != nil {
		return ""
	}
	defer dom.Free()
	instXml := libvirtc.NewDomainXMLFromDom(dom, true)
	if instXml == nil {
		return ""
	}
	if _, port := instXml.GraphicsAddr(typ); port != "" {
		return hyper.Address + ":" + port
	}
	return ""
}

func (w WebSocket) GetTarget(r *http.Request) string {
	id := api.GetQueryOne(r, "id")
	if id == "" {
		return ""
	}
	format := api.GetQueryOne(r, "type")
	if format == "" {
		format = "vnc"
	}
	name := api.GetQueryOne(r, "node")
	if name == "" {
		return w.GetLocal(id, format)
	}
	return w.GetRemote(id, name, format)
}

func (w WebSocket) Handle(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	target := w.GetTarget(ws.Request())
	if target == "" {
		libstar.Error("WebSocket.Handle target not found.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("WebSocket.Handle dial %s", err)
		return
	}
	defer conn.Close()
	libstar.Info("WebSocket.Handle request by %s", ws.RemoteAddr())
	libstar.Info("WebSocket.Handle connect to %s", conn.RemoteAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Warn("WebSocket.Handle copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Warn("WebSocket.Handle copy from target %s", err)
		}
	}()
	wait.Wait()
}

type TcpSocket struct {
}

func (t TcpSocket) Router(router *mux.Router) {
	router.Handle("/ext/tcpsocket", websocket.Handler(t.Handle))
}

func (t TcpSocket) Local(host string, ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	r := ws.Request()
	target := api.GetQueryOne(r, "target")
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("TcpSocket.Local dial %s", err)
		return
	}
	defer conn.Close()
	user, _ := api.GetUser(r)
	libstar.Info("TcpSocket.Local with %s", user.Name)
	libstar.Info("TcpSocket.Local request by %s", ws.RemoteAddr())
	libstar.Info("TcpSocket.Local connect to %s", conn.RemoteAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Warn("TcpSocket.Local copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Warn("TcpSocket.Local copy from target %s", err)
		}
	}()
	wait.Wait()
}

func (t TcpSocket) Remote(host string, ws *websocket.Conn) {
	r := ws.Request()
	node := service.SERVICE.Zone.Get(host)
	if node == nil {
		libstar.Error("TcpSocket.Remote host not found: %s", host)
		return
	}
	query := r.URL.Query()
	query.Set("host", "")
	r.URL.RawQuery = query.Encode()
	pri := libstar.ProxyWs{
		Proxy: libstar.Proxy{
			Server: node.Url,
			Auth: libstar.Auth{
				Type:     "basic",
				Username: node.Username,
				Password: node.Password,
			},
		},
	}
	pri.Initialize()
	pri.Socket(ws)
}

func (t TcpSocket) Handle(ws *websocket.Conn) {
	r := ws.Request()
	host := api.GetQueryOne(r, "host")
	if host == "" {
		t.Local("", ws)
	} else {
		t.Remote(host, ws)
	}
}
