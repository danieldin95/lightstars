package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/compute/libvirt"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
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
	return
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
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r,
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
	libstar.Debug("Server.IsAuth %s:%s", user, pass)
	if !ok || pass != h.adminToken || user != "admin" {
		w.Header().Set("WWW-Authenticate", "Basic")
		http.Error(w, "Authorization Required.", http.StatusUnauthorized)
		return false
	}
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
	router := h.Router()
	staticFile := http.StripPrefix("/static/", http.FileServer(http.Dir(h.pubDir)))

	router.HandleFunc("/", h.HandleIndex)
	router.PathPrefix("/static/").Handler(staticFile)
	router.Handle("/websockify", websocket.Handler(h.HandleWebSockify))
	router.HandleFunc("/api/instance/{id}", h.GetInstance).Methods("GET")
	router.HandleFunc("/api/instance/{id}", h.AddInstance).Methods("POST")
	router.HandleFunc("/api/instance/{id}", h.ModInstance).Methods("PUT")
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

func (h *Server) ParseFiles(w http.ResponseWriter, name string, data interface{}) error {
	file := path.Base(name)
	tmpl, err := template.New(file).Funcs(template.FuncMap{
		"prettyBytes":  libstar.PrettyBytes,
		"prettyKBytes": libstar.PrettyKBytes,
		"prettySecs":   libstar.PrettySecs,
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
	index := IndexSchema{
		Instances: make([]InstanceSchema, 0, 32),
	}

	hyper, err := libvirtdriver.GetHyper("")
	if err != nil {
		libstar.Error("Server.HandleIndex %s", err)
		return
	}
	index.Version = NewVersionSchema()
	if domains, err := hyper.ListAllDomains(); err == nil {
		for _, dom := range domains {
			instance := NewInstanceSchema(dom)
			index.Instances = append(index.Instances, instance)
		}
	}
	file := h.GetFile("/index.html")
	if err := h.ParseFiles(w, file, index); err != nil {
		libstar.Error("Server.HandleIndex %s", err)
	}
}

func (h *Server) GetTarget(req *http.Request) string {
	var id string

	query := req.URL.Query()
	if tgt, ok := query["target"]; ok {
		return tgt[0]
	}
	ids, ok := query["instance"]
	if ok {
		id = ids[0]
	}
	libstar.Info("Server.GetTarget %s", id)
	hyper, err := libvirtdriver.GetHyper("")
	if err != nil {
		libstar.Error("Server.HandleIndex %s", err)
		return ""
	}
	dom, err := hyper.LookupDomainByUUIDString(id)
	if err != nil {
		return ""
	}
	instXml := libvirtdriver.NewDomainXMLFromDom(dom, true)
	if instXml == nil {
		return ""
	}
	if _, port := instXml.VNCDisplay(); port != "" {
		return hyper.Address + ":" + port
	}
	return ""
}

func (h *Server) HandleWebSockify(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	target := h.GetTarget(ws.Request())
	if target == "" {
		libstar.Error("Server.HandleWebSockify target not found.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("Server.HandleWebSockify dial %s", err)
		return
	}
	defer conn.Close()

	libstar.Info("Server.HandleWebSockify request from %s", ws.RemoteAddr())
	libstar.Info("Server.HandleWebSockify connection to %s", conn.LocalAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)

	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Error("Server.HandleWebSockify copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Error("Server.HandleWebSockify copy from target %s", err)
		}
	}()
	wait.Wait()
}

func (h *Server) GetArg(r *http.Request, name string) (string, bool) {
	vars := mux.Vars(r)
	value, ok := vars[name]
	return value, ok
}

func (h *Server) GetInstance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := h.GetArg(r, "id")

	dom, err := libvirtdriver.LookupDomainByUUIDString("", uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.ResponseJson(w, libvirtdriver.NewDomainXMLFromDom(dom, true))
}

func (h *Server) AddInstance(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//user := &UserSchema{}
	//if err := json.Unmarshal([]byte(body), user); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//h.ResponseMsg(w, 0, "")
}

func (h *Server) ModInstance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := h.GetArg(r, "id")

	hyper, err := libvirtdriver.GetHyper("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conf := &InstanceConfSchema{}
	if err := json.Unmarshal([]byte(body), conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch conf.Action {
	case "start":
		xmlData, err := dom.GetXMLDesc(0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := dom.Undefine(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if dom, err := hyper.DomainDefineXML(xmlData); err == nil {
			if err := dom.Create(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "shutdown":
		xmlData, err := dom.GetXMLDesc(0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := dom.Shutdown(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if _, err := hyper.DomainDefineXML(xmlData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "suspend":
		if err := dom.Suspend(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "resume":
		if err := dom.Resume(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	h.ResponseMsg(w, 0, "")
}