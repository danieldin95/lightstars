package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"sync"
	"text/template"
	"time"
)

type Server struct {
	listen     string
	server     *http.Server
	crtFile    string
	keyFile    string
	pubDir     string
	router     *mux.Router
}

func NewServer(listen string) (h *Server) {
	h = &Server{
		listen: listen,
	}

	return
}

func (h *Server) Initialize() {
	r := h.Router()
	if h.server == nil {
		h.server = &http.Server{
			Addr:         h.listen,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r,
		}
	}

	h.LoadRouter()
}

func (h *Server) IsAuth(w http.ResponseWriter, r *http.Request) bool {
	//token, pass, ok := r.BasicAuth()
	//libstar.Debug("Server.IsAuth token: %s, pass: %s", token, pass)
	//
	//if len(r.URL.Path) < 4 || r.URL.Path[:4] != "/api" {
	//	return true
	//}
	//
	//if !ok || token != h.adminToke { //
	//	w.Header().Set("WWW-Authenticate", "Basic")
	//	http.Error(w, "Authorization Required.", http.StatusUnauthorized)
	//	return false
	//}
	return true
}

func (h *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.IsAuth(w, r) {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic")
			http.Error(w, "Authorization Required.", http.StatusUnauthorized)
		}
	})
}

func (h *Server) Router() *mux.Router {
	if h.router == nil {
		h.router = mux.NewRouter()
		h.router.Use(h.Middleware)
	}

	return h.router
}

func (h *Server) LoadRouter() {
	h.Router().HandleFunc("/", h.HandleIndex)
	h.Router().HandleFunc("/favicon.ico", h.PubFile)
	h.Router().Handle("/websockify", websocket.Handler(h.HandleWebsockify))
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

func (h *Server) ResponseJson(w http.ResponseWriter, v interface{}) {
	str, err := json.Marshal(v)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(str)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Server) ResponseMsg(w http.ResponseWriter, code int, message string) {
	ret := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	}
	h.ResponseJson(w, ret)
}

func (h *Server) GetFile(name string) string {
	return fmt.Sprintf("%s%s", h.pubDir, name)
}

func (h *Server) PubFile(w http.ResponseWriter, r *http.Request) {
	realpath := h.GetFile(r.URL.Path)
	contents, err := ioutil.ReadFile(realpath)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	fmt.Fprintf(w, "%s\n", contents)
}

func (h *Server) ParseFiles(w http.ResponseWriter, name string, data interface{}) error {
	file := path.Base(name)
	tmpl, err := template.New(file).Funcs(template.FuncMap{
	}).ParseFiles(name)
	if err != nil {
		fmt.Fprintf(w, "template.ParseFiles %s", err)
		return err
	}
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Fprintf(w, "template.ParseFiles %s", err)
		return err
	}
	return nil
}

func (h *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Index\n"))
}

func (h *Server) GetTarget(id string) string {
	return "192.168.4.249:5900"
}
func (h *Server) HandleWebsockify(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	conn, err := net.Dial("tcp", h.GetTarget(""))
	if err != nil {
		libstar.Error("Server.HandleWebsockify dial %s", err)
		return
	}
	defer conn.Close()

	libstar.Info("Server.HandleWebsockify request from %s", ws.RemoteAddr())
	libstar.Info("Server.HandleWebsockify connection to %s", conn.LocalAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)

	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Error("Server.HandleWebsockify copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Error("Server.HandleWebsockify copy from target %s", err)
		}
	}()
	wait.Wait()
}