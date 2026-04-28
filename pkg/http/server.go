package http

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/danieldin95/lightstar/pkg/http/api"
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/service"
	"github.com/danieldin95/lightstar/pkg/storage"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type historyResponseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *historyResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *historyResponseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	if w.body.Len() < 8192 {
		remain := 8192 - w.body.Len()
		if remain > len(p) {
			remain = len(p)
		}
		_, _ = w.body.Write(p[:remain])
	}
	return w.ResponseWriter.Write(p)
}

func (w *historyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

func (w *historyResponseWriter) Flush() {
	if fl, ok := w.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

func (w *historyResponseWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := w.ResponseWriter.(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return http.ErrNotSupported
}

type Server struct {
	listen     string
	server     *http.Server
	crtFile    string
	keyFile    string
	pubDir     string
	router     *mux.Router
	adminToken string
	adminFile  string
}

func NewServer(listen, staticDir, authFile string) (h *Server) {
	h = &Server{
		listen:    listen,
		pubDir:    staticDir,
		adminFile: authFile,
	}
	api.SetStatic(h.pubDir)
	return
}

func (h *Server) Router() *mux.Router {
	if h.router != nil {
		return h.router
	}
	h.router = mux.NewRouter()
	return h.router
}

func (h *Server) LoadRouter() {
	router := h.Router()
	router.NotFoundHandler = http.HandlerFunc(h.Handle404)
	router.Use(h.Middleware)

	// static files
	Static{}.Router(router)
	Download{}.Router(router)
	// proxy websocket
	WsGraphics{}.Router(router)
	WsTcp{}.Router(router)
	WsProxy{}.Router(router)
	// zone router
	Host{}.Router(router)
	// ui router
	UI{}.Router(router)
	Login{}.Router(router)
	// api router
	api.Upload{}.Router(router)
	api.Hyper{}.Router(router)
	api.ISO{}.Router(router)
	api.Volume{}.Router(router)
	api.Bridger{}.Router(router)
	api.DataStore{}.Router(router)
	api.Network{}.Router(router)
	api.Instance{}.Router(router)
	api.Disk{}.Router(router)
	api.Interface{}.Router(router)
	api.Processor{}.Router(router)
	api.Memory{}.Router(router)
	api.Graphics{}.Router(router)
	api.Zone{}.Router(router)
	api.ProxyTcp{}.Router(router)
	api.DHCPLease{}.Router(router)
	api.Volume{}.Router(router)
	api.User{}.Router(router)
	api.Snapshot{}.Router(router)
	api.History{}.Router(router)
}

func (h *Server) SetCert(keyFile, crtFile string) {
	h.crtFile = crtFile
	h.keyFile = keyFile
}

func (h *Server) Initialize() {
	r := h.Router()
	if h.server == nil {
		h.server = &http.Server{
			Addr:    h.listen,
			Handler: r,
		}
	}
	path := storage.PATH.RootXML()
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			libstar.Warn("Server.Initialize %s", err)
		}
	}
	h.LoadRouter()
}

func (h *Server) Filter(r *http.Request) bool {
	path := r.URL.Path
	libstar.Debug("Server.Filter %v", r.Header)
	if strings.HasPrefix(path, "/static/sshy") {
		return false
	}
	if path == "/" || strings.HasPrefix(path, "/static") ||
		strings.HasPrefix(path, "/favicon.ico") ||
		strings.HasPrefix(path, "/ui/login") ||
		strings.HasPrefix(path, "/ui/console") ||
		strings.HasPrefix(path, "/ext/ws/graphics") {
		return true
	}
	return false
}

func (h *Server) IsAuth(w http.ResponseWriter, r *http.Request) bool {
	// not need to auth.
	if h.Filter(r) {
		return true
	}
	name, pass, _ := api.GetAuth(r)
	libstar.Print("Server.IsAuth %s:%s", name, pass)
	// auth by password and name
	user, ok := service.SERVICE.Users.Get(name)
	if !ok || user.Password != pass {
		return false
	}
	return true
}

func (h *Server) LogRequest(r *http.Request) {
	path := r.URL.Path
	if q, _ := url.QueryUnescape(r.URL.RawQuery); q != "" {
		path += "?" + q
	}
	if strings.HasPrefix(r.URL.Path, "/static") ||
		strings.HasSuffix(r.URL.Path, ".ico") ||
		strings.HasSuffix(r.URL.Path, ".png") ||
		strings.HasSuffix(r.URL.Path, ".gif") {
		libstar.Debug("Server.Middleware %s %s %s", r.RemoteAddr, r.Method, path)
	} else if r.Method == "GET" {
		libstar.Debug("Server.Middleware %s %s %s", r.RemoteAddr, r.Method, path)
	} else {
		libstar.Info("Server.Middleware %s %s %s", r.RemoteAddr, r.Method, path)
	}
}

func (h *Server) History(user schema.User, r *http.Request, status int, body string) {
	if strings.HasPrefix(r.URL.Path, "/ui") {
		return
	}
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
		method := strings.ToUpper(strings.TrimSpace(r.Method))
		if method == "" {
			method = "UNKNOWN"
		}
		path := strings.TrimSpace(r.URL.Path)
		if path == "" {
			path = strings.TrimSpace(r.RequestURI)
		}
		if path == "" {
			path = "-"
		}
		result := "success"
		trimBody := strings.TrimSpace(body)
		if trimBody != "" {
			var resp struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal([]byte(trimBody), &resp); err == nil {
				if strings.TrimSpace(resp.Message) != "" {
					result = strings.TrimSpace(resp.Message)
				}
			}
		}
		if status >= 400 {
			result = trimBody
			if result == "" {
				result = http.StatusText(status)
			}
		}
		his := &schema.History{
			User:   user.Name,
			Date:   time.Now().Format(time.RFC3339),
			Method: method,
			Url:    path,
			Client: r.RemoteAddr,
			Result: result,
		}
		service.SERVICE.History.Add(his)
	}
}

func (h *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.LogRequest(r)
		if h.IsAuth(w, r) {
			user, _ := api.GetUser(r)
			if user.Type == "admin" || service.SERVICE.Permission.Has(r) {
				api.UpdateCookie(w, r, user)
				rw := &historyResponseWriter{ResponseWriter: w}
				next.ServeHTTP(rw, r)
				h.History(user, r, rw.status, rw.body.String())
			} else {
				http.Error(w, "Request not allowed", http.StatusForbidden)
			}
		} else if strings.HasPrefix(r.URL.Path, "/ui") {
			http.Redirect(w, r, "/ui/login", http.StatusMovedPermanently)
		} else {
			//w.Header().Set("WWW-Authenticate", "Basic")
			http.Error(w, "Authorization Required", http.StatusUnauthorized)
		}
	})
}

func (h *Server) Start() {
	h.Initialize()
	if h.keyFile == "" || h.crtFile == "" {
		libstar.Info("Server.Start http://%s", h.listen)
		if err := h.server.ListenAndServe(); err != nil {
			libstar.Error("Server.Start on %s: %s", h.listen, err)
			return
		}
	} else {
		libstar.Info("Server.Start https://%s", h.listen)
		if err := h.server.ListenAndServeTLS(h.crtFile, h.keyFile); err != nil {
			libstar.Error("Server.Start on %s: %s", h.listen, err)
			return
		}
	}
}

func (h *Server) Shutdown() {
	libstar.Info("Server.Shutdown %s", h.listen)
	if err := h.server.Shutdown(context.Background()); err != nil {
		libstar.Error("Server.Shutdown %v", err)
	}
}

func (h *Server) Handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	file := api.GetFile("ui/404.html")
	if err := api.ParseFiles(w, file, nil); err != nil {
		libstar.Error("Server.Handle404 %s", err)
	}
}
