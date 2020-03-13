package http

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
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
	router.Handle("/ext/websocket", websocket.Handler(w.Socket))
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
	libstar.Debug("UI.GetTarget %s://%s", format, id)
	hyper, err := libvirtc.GetHyper()
	if err != nil {
		libstar.Error("UI.GetTarget %s", err)
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
	if _, port := instXml.GraphicsAddr(format); port != "" {
		return hyper.Address + ":" + port
	}
	return ""
}

func (w WebSocket) Socket(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	target := w.GetTarget(ws.Request())
	if target == "" {
		libstar.Error("UI.Socket target not found.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("UI.Socket dial %s", err)
		return
	}
	defer conn.Close()
	libstar.Info("UI.Socket request by %s", ws.RemoteAddr())
	libstar.Info("UI.Socket connect to %s", conn.RemoteAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Error("UI.Socket copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Error("UI.Socket copy from target %s", err)
		}
	}()
	wait.Wait()
}
