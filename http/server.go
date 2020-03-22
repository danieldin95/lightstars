package http

import (
	"context"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/service"
	"github.com/danieldin95/lightstar/storage"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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
	WebSocket{}.Router(router)
	TcpSocket{}.Router(router)
	// zone router
	Host{}.Router(router)
	// ui router
	UI{}.Router(router)
	Login{}.Router(router)
	// api router
	api.Upload{}.Router(router)
	api.Hyper{}.Router(router)
	api.ISO{}.Router(router)
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

func (h *Server) IsAuth(w http.ResponseWriter, r *http.Request) bool {
	// not need to auth.
	if r.URL.Path == "/" ||
		strings.HasPrefix(r.URL.Path, "/ext/webs") ||
		strings.HasPrefix(r.URL.Path, "/static") ||
		strings.HasPrefix(r.URL.Path, "/favicon.ico") ||
		strings.HasPrefix(r.URL.Path, "/ui/login") ||
		strings.HasPrefix(r.URL.Path, "/ui/console") {
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
	if strings.HasPrefix(r.URL.Path, "/static") ||
		strings.HasSuffix(r.URL.Path, ".ico") ||
		strings.HasSuffix(r.URL.Path, ".png") ||
		strings.HasSuffix(r.URL.Path, ".gif") {
		return
	}
	path := r.URL.Path
	if q, _ := url.QueryUnescape(r.URL.RawQuery); q != "" {
		path += "?" + q
	}
	libstar.Info("Server.Middleware %s %s %s", r.RemoteAddr, r.Method, path)
}

func (h *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.LogRequest(r)
		if h.IsAuth(w, r) {
			user, _ := api.GetUser(r)
			if user.Type == "admin" || service.SERVICE.Permission.Has(r) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Request not allowed", http.StatusForbidden)
			}
		} else if strings.HasPrefix(r.URL.Path, "/ui") {
			http.Redirect(w, r, "/ui/login", http.StatusMovedPermanently)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic")
			http.Error(w, "Authorization Required", http.StatusUnauthorized)
		}
	})
}

func (h *Server) Start() error {
	h.Initialize()
	if h.keyFile == "" || h.crtFile == "" {
		libstar.Info("Server.Start http://%s", h.listen)
		if err := h.server.ListenAndServe(); err != nil {
			libstar.Error("Server.Start on %s: %s", h.listen, err)
			return err
		}
	} else {
		libstar.Info("Server.Start https://%s", h.listen)
		if err := h.server.ListenAndServeTLS(h.crtFile, h.keyFile); err != nil {
			libstar.Error("Server.Start on %s: %s", h.listen, err)
			return err
		}
	}
	return nil
}

func (h *Server) Shutdown() {
	libstar.Info("Server.Shutdown %s", h.listen)
	if err := h.server.Shutdown(context.Background()); err != nil {
		libstar.Error("Server.Shutdown %v", err)
	}
}

func (h *Server) Handle404(w http.ResponseWriter, r *http.Request) {
	file := api.GetFile("ui/404.html")
	if err := api.ParseFiles(w, file, nil); err != nil {
		libstar.Error("Server.Handle404 %s", err)
	}
}
