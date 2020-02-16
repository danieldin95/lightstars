package vswitch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"
	"time"
)

type Http struct {
	listen     string
	server     *http.Server
	crtFile    string
	keyFile    string
	pubDir     string
	router     *mux.Router
}

func NewHttp() (h *Http) {
	h = &Http{
	}

	return
}

func (h *Http) Initialize() {
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

func (h *Http) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.IsAuth(w, r) {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic")
			http.Error(w, "Authorization Required.", http.StatusUnauthorized)
		}
	})
}

func (h *Http) Router() *mux.Router {
	if h.router == nil {
		h.router = mux.NewRouter()
		h.router.Use(h.Middleware)
	}

	return h.router
}

func (h *Http) LoadRouter() {
	h.Router().HandleFunc("/", h.IndexHtml)
	h.Router().HandleFunc("/favicon.ico", h.PubFile)
}

func (h *Http) Start() error {
	h.Initialize()

	libstar.Info("Http.Start %s", h.listen)

	if h.keyFile == "" || h.crtFile == "" {
		if err := h.server.ListenAndServe(); err != nil {
			libstar.Error("Http.Start on %s: %s", h.listen, err)
			return err
		}
	} else {
		if err := h.server.ListenAndServeTLS(h.crtFile, h.keyFile); err != nil {
			libstar.Error("Http.Start on %s: %s", h.listen, err)
			return err
		}
	}
	return nil
}

func (h *Http) Shutdown() {
	libstar.Info("Http.Shutdown %s", h.listen)
	if err := h.server.Shutdown(context.Background()); err != nil {
		libstar.Error("Http.Shutdown: %v", err)
	}
}

func (h *Http) IsAuth(w http.ResponseWriter, r *http.Request) bool {
	token, pass, ok := r.BasicAuth()
	libstar.Debug("Http.IsAuth token: %s, pass: %s", token, pass)

	if len(r.URL.Path) < 4 || r.URL.Path[:4] != "/api" {
		return true
	}

	if !ok { // || token != h.adminToke
		w.Header().Set("WWW-Authenticate", "Basic")
		http.Error(w, "Authorization Required.", http.StatusUnauthorized)
		return false
	}

	return true
}

func (h *Http) ResponseJson(w http.ResponseWriter, v interface{}) {
	str, err := json.Marshal(v)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(str)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Http) ResponseMsg(w http.ResponseWriter, code int, message string) {
	ret := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	}
	h.ResponseJson(w, ret)
}

func (h *Http) getFile(name string) string {
	return fmt.Sprintf("%s%s", h.pubDir, name)
}

func (h *Http) PubFile(w http.ResponseWriter, r *http.Request) {
	realpath := h.getFile(r.URL.Path)
	contents, err := ioutil.ReadFile(realpath)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}

	fmt.Fprintf(w, "%s\n", contents)
}

func (h *Http) ParseFiles(w http.ResponseWriter, name string, data interface{}) error {
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

func (h *Http) IndexHtml(w http.ResponseWriter, r *http.Request) {
}