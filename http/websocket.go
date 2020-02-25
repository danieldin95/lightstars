package http

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/http"
	"sync"
)

type WebSocket struct {
}

func (w WebSocket) Router(router *mux.Router) {
	router.Handle("/websockify", websocket.Handler(w.Sockify))
}

func (w WebSocket) GetTarget(r *http.Request) string {
	if t := api.GetQueryOne(r, "target"); t != "" {
		return t
	}
	id := api.GetQueryOne(r, "instance")
	if id == "" {
		return ""
	}

	libstar.Debug("UI.GetTarget %s", id)
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
	if _, port := instXml.VNCDisplay(); port != "" {
		return hyper.Address + ":" + port
	}
	return ""
}

func (w WebSocket) Sockify(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	target := w.GetTarget(ws.Request())
	if target == "" {
		libstar.Error("UI.Sockify target not found.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("UI.Sockify dial %s", err)
		return
	}
	defer conn.Close()
	libstar.Info("UI.Sockify request by %s", ws.RemoteAddr())
	libstar.Info("UI.Sockify connect to %s", conn.RemoteAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Error("UI.Sockify copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Error("UI.Sockify copy from target %s", err)
		}
	}()
	wait.Wait()
}
