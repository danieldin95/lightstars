package http

import (
	"encoding/base64"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/http/api"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/service"
	"github.com/danieldin95/lightstar/src/storage"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/http"
)

type Download struct {
}

func (down Download) Router(router *mux.Router) {
	dir := http.Dir(storage.PATH.Root())
	files := http.StripPrefix("/ext/files/", FileServer(dir))
	router.PathPrefix("/ext/files/").Handler(files)
}

func proxy(src, dst net.Conn) {
	libstar.Debug("proxy by %s", src.RemoteAddr())
	libstar.Debug("proxy to %s", dst.RemoteAddr())
	wait := libstar.NewWaitOne(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(dst, src); err != nil {
			libstar.Debug("proxy from %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(src, dst); err != nil {
			libstar.Debug("proxy from %s", err)
		}
	}()
	wait.Wait()
	libstar.Debug("proxy %s exit", src.RemoteAddr())
}

type WsGraphics struct {
}

func (w WsGraphics) Router(router *mux.Router) {
	router.Handle("/ext/ws/graphics", websocket.Handler(w.Handle))
}

func (w WsGraphics) GetRemote(id, name, typ string) string {
	libstar.Debug("WsGraphics.GetRemote %s://%s@%s", typ, id, name)
	node := service.SERVICE.Zone.Get(name)
	if node == nil {
		libstar.Error("WsGraphics.GetRemote %s", name)
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
		libstar.Error("WsGraphics.GetRemote %s", err)
		return ""
	}
	defer client.Close()
	inst := schema.Instance{}
	if err := libstar.GetJSON(resp.Body, &inst); err != nil {
		libstar.Error("WsGraphics.GetRemote %s", name)
		return ""
	}
	port := ""
	for _, g := range inst.Graphics {
		if typ == g.Type {
			port = g.Port
		}
	}
	if port == "" || port == "-1" {
		return ""
	}
	return host + ":" + port
}

func (w WsGraphics) GetLocal(id, typ string) string {
	libstar.Debug("WsGraphics.GetLocal %s://%s", typ, id)
	hyper, err := libvirtc.GetHyper()
	if err != nil {
		libstar.Error("WsGraphics.GetLocal %s", err)
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
	if _, port := instXml.GraphicsAddr(typ); port == "" || port == "-1" {
		return ""
	} else {
		return hyper.Address + ":" + port
	}
}

func (w WsGraphics) GetTarget(r *http.Request) string {
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

func (w WsGraphics) Handle(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame
	target := w.GetTarget(ws.Request())
	if target == "" {
		libstar.Debug("WsGraphics.Handle target not found.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("WsGraphics.Handle %s", err)
		return
	}
	defer conn.Close()
	proxy(ws, conn)
}

type WsTcp struct {
}

func (t WsTcp) Router(router *mux.Router) {
	router.Handle("/ext/ws/tcp", websocket.Handler(t.Handle))
}

func (t WsTcp) Local(host string, ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame
	r := ws.Request()
	target := api.GetQueryOne(r, "target")
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("WsTcp.Local %s", err)
		return
	}
	defer conn.Close()
	user, _ := api.GetUser(r)
	libstar.Debug("WsTcp.Local with %s", user.Name)
	proxy(ws, conn)
}

func (t WsTcp) Remote(host string, ws *websocket.Conn) {
	r := ws.Request()
	node := service.SERVICE.Zone.Get(host)
	if node == nil {
		libstar.Error("WsTcp.Remote host not found: %s", host)
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

func (t WsTcp) Handle(ws *websocket.Conn) {
	r := ws.Request()
	host := api.GetQueryOne(r, "host")
	if host == "" {
		t.Local("", ws)
	} else {
		t.Remote(host, ws)
	}
}

type WsProxy struct {
}

func (w WsProxy) Router(router *mux.Router) {
	router.Handle("/ext/ws/proxy/{target}", websocket.Handler(w.Handle))
}

func (w WsProxy) Handle(ws *websocket.Conn) {
	defer ws.Close()
	target, _ := api.GetArg(ws.Request(), "target")
	if target == "" {
		libstar.Error("WsProxy.Handle target notFound.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("WsProxy.Handle %s", err)
		return
	}
	defer conn.Close()
	libstar.Debug("WsProxy.Handle %s", target)
	wait := libstar.NewWaitOne(2)
	go func() {
		defer wait.Done()
		for {
			var data []byte
			if err := websocket.Message.Receive(ws, &data); err != nil {
				libstar.Warn("WsProxy.Handle.ws.Receive %s", err)
				break
			}
			message := make([]byte, len(data))
			n, err := base64.StdEncoding.Decode(message, data)
			if err != nil {
				libstar.Warn("WsProxy.Handle.base64.Decode %s", err)
				break
			}
			if _, err := conn.Write(message[:n]); err != nil {
				libstar.Warn("WsProxy.Handle.conn.Write %s", err)
				break
			}
		}
	}()
	go func() {
		defer wait.Done()
		for {
			message := make([]byte, 4096)
			n, err := conn.Read(message)
			if err != nil {
				libstar.Warn("WsProxy.Handle.conn.Read %s", err)
				break
			}
			data := base64.StdEncoding.EncodeToString(message[:n])
			if err := websocket.Message.Send(ws, data); err != nil {
				libstar.Warn("WsProxy.Handle.ws.Send %s", err)
				break
			}
		}
	}()
	wait.Wait()
}
