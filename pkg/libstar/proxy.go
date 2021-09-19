package libstar

import (
	"crypto/tls"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"net/url"
)

type Auth struct {
	Type     string
	Username string
	Password string
}

type Proxy struct {
	Auth      Auth
	Prefix    string
	Server    string
	Url       *url.URL
	TlsConfig *tls.Config
}

func (pri *Proxy) Initialize() {
	pri.Url, _ = url.Parse(pri.Server)
	pri.TlsConfig = &tls.Config{InsecureSkipVerify: true}
}

func (pri *Proxy) GetPath(req *http.Request) string {
	path := req.URL.Path
	size := len(pri.Prefix)
	if size > 0 && len(path) >= size {
		path = path[size:]
	}
	return path + "?" + req.URL.RawQuery
}

type ProxyUrl struct {
	Proxy
	Transport *http.Transport
	Filter    func(*http.Response, http.ResponseWriter, interface{}) bool
	Data      interface{}
}

func (pri *ProxyUrl) Initialize() {
	pri.Proxy.Initialize()
	pri.Transport = &http.Transport{
		TLSClientConfig: pri.TlsConfig,
	}
}

func (pri *ProxyUrl) ServeHttp(w http.ResponseWriter, r *http.Request) {
	p, err := pri.Transport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer pri.Transport.CloseIdleConnections()
	defer p.Body.Close()
	for key, value := range p.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	status := p.StatusCode
	w.WriteHeader(status)
	filtered := false
	if pri.Filter != nil {
		filtered = pri.Filter(p, w, pri.Data)
	}
	if !filtered {
		_, _ = io.Copy(w, p.Body)
	}
}

func (pri *ProxyUrl) Handler(w http.ResponseWriter, req *http.Request) {
	urlPath := pri.Server
	urlPath += pri.GetPath(req)
	Debug("ProxyUrl.Handler %s %s to %s", req.Method, req.URL.Path, urlPath)
	outReq, _ := http.NewRequest(req.Method, urlPath, req.Body)
	for key, value := range req.Header {
		for _, v := range value {
			outReq.Header.Add(key, v)
		}
		Print("ProxyUrl.Handler %s: %s", key, value)
	}
	Debug("ProxyUrl.Handler %s", pri.Auth)
	if pri.Auth.Type == "basic" {
		outReq.SetBasicAuth(pri.Auth.Username, pri.Auth.Password)
	}
	pri.ServeHttp(w, outReq)
}

type ProxyWs struct {
	Proxy
}

func (pri *ProxyWs) Dial(url_, protocol, origin string) (ws *websocket.Conn, err error) {
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		return nil, err
	}
	if protocol != "" {
		config.Protocol = []string{protocol}
	}
	config.TlsConfig = pri.TlsConfig
	if pri.Auth.Type == "basic" {
		config.Header = http.Header{
			"Authorization": {BasicAuth(pri.Auth.Username, pri.Auth.Password)},
		}
	}
	return websocket.DialConfig(config)
}

func (pri *ProxyWs) Socket(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	req := ws.Request()
	schema := "ws"
	if pri.Url.Scheme == "https" {
		schema = "wss"
	}
	target := schema + "://" + pri.Url.Host + pri.GetPath(req)
	conn, err := pri.Dial(target, "", req.URL.RequestURI())
	if err != nil {
		Error("ProxyWs.Socket %s", err)
		return
	}
	defer conn.Close()
	Info("ProxyWs.Socket by %s", ws.RemoteAddr())
	Info("ProxyWs.Socket to %s", conn.RemoteAddr())

	wait := NewWaitOne(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			Warn("ProxyWs.Socket from ws %v", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			Warn("ProxyWs.Socket from target %v", err)
		}
	}()
	wait.Wait()
	Warn("ProxyWs.Socket %s exit", ws.RemoteAddr())
}
