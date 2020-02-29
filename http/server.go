package http

import (
	"context"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
	if h.adminToken == "" {
		h.LoadToken()
	}
	if h.adminToken == "" {
		h.adminToken = libstar.GenToken(64)
	}
	h.SaveToken()
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
	h.router.NotFoundHandler = http.HandlerFunc(h.Handle404)
	h.router.Use(h.Middleware)

	// static files
	Static{}.Router(router)
	// proxy websocket
	WebSocket{}.Router(router)
	// ui router
	UI{}.Router(router)
	// api router
	api.ISO{}.Router(router)
	api.Bridger{}.Router(router)
	api.DataStore{}.Router(router)
	api.Network{}.Router(router)
	api.Disk{}.Router(router)
	api.Interface{}.Router(router)
	api.Instance{}.Router(router)
}

func (h *Server) SetCert(keyFile, crtFile string) {
	h.crtFile = crtFile
	h.keyFile = keyFile
}

func (h *Server) Initialize() {
	r := h.Router()
	if h.server == nil {
		h.server = &http.Server{
			Addr:         h.listen,
			WriteTimeout: time.Second * 60 * 5,
			ReadTimeout:  time.Second * 60 * 5,
			IdleTimeout:  time.Second * 60 * 30,
			Handler:      r,
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

func (h *Server) SaveToken() error {
	libstar.Info("Server.SaveToken: AdminToken: %s", h.adminToken)
	f, err := os.OpenFile(h.adminFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		libstar.Error("Server.SaveToken: %s", err)
		return err
	}
	if _, err := f.Write([]byte(h.adminToken)); err != nil {
		libstar.Error("Server.SaveToken: %s", err)
		return err
	}
	return nil
}

func (h *Server) LoadToken() error {
	if _, err := os.Stat(h.adminFile); os.IsNotExist(err) {
		libstar.Info("Server.LoadToken: file:%s does not exist", h.adminFile)
		return nil
	}
	contents, err := ioutil.ReadFile(h.adminFile)
	if err != nil {
		libstar.Error("Server.LoadToken: file:%s %s", h.adminFile, err)
		return err

	}
	h.adminToken = string(contents)
	return nil
}

func (h *Server) IsAuth(w http.ResponseWriter, r *http.Request) bool {
	user, pass, ok := r.BasicAuth()
	libstar.Print("Server.IsAuth %s:%s", user, pass)

	if strings.HasPrefix(r.URL.Path, "/static") {
		return true
	} else if r.URL.Path == "/ui/console" || r.URL.Path == "/websockify" {
		return true
	}
	if !ok || pass != h.adminToken || user != "admin" {
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
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic")
			http.Error(w, "Authorization Required", http.StatusUnauthorized)
		}
	})
}

func (h *Server) Start() error {
	h.Initialize()
	libstar.Info("Server.Start %s", h.listen)
	if h.keyFile == "" || h.crtFile == "" {
		if err := h.server.ListenAndServe(); err != nil {
			libstar.Error("Server.Start on %s: %s", h.listen, err)
			return err
		}
	} else {
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
