package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/compute/libvirt"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage/qemu"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
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
	router.HandleFunc("/api/instance", h.AddInstance).Methods("POST")
	router.HandleFunc("/api/instance/{id}", h.GetInstance).Methods("GET")
	router.HandleFunc("/api/instance/{id}", h.ModInstance).Methods("PUT")
	router.HandleFunc("/api/iso", h.GetISO).Methods("GET")
	router.HandleFunc("/api/bridge", h.GetBridge).Methods("GET")
	router.HandleFunc("/api/datastore", h.GetDataStore).Methods("GET")
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

func (h *Server) ResponseXml(w http.ResponseWriter, v string) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(v))
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

	hyper, err := libvirtdriver.GetHyper()
	if err != nil {
		libstar.Error("Server.HandleIndex %s", err)
		return
	}
	index.Version = NewVersionSchema()
	index.Hyper = NewHyperSchema()
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

func (h *Server) GetQueryOne(req *http.Request, name string) string {
	query := req.URL.Query()
	if values, ok := query[name]; ok {
		return values[0]
	}
	return ""
}

func (h *Server) GetTarget(req *http.Request) string {
	if t := h.GetQueryOne(req, "target"); t != "" {
		return t
	}
	id := h.GetQueryOne(req, "instance")
	if id == "" {
		return ""
	}

	libstar.Info("Server.GetTarget %s", id)
	hyper, err := libvirtdriver.GetHyper()
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

	dom, err := libvirtdriver.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	format := h.GetQueryOne(r, "format")
	if format == "xml" {
		xmlDesc, err := dom.GetXMLDesc(false)
		if err == nil {
			h.ResponseXml(w, xmlDesc)
		} else {
			h.ResponseXml(w, "<error>"+err.Error()+"</error>")
		}
	} else {
		h.ResponseJson(w, libvirtdriver.NewDomainXMLFromDom(dom, true))
	}
}

func (h *Server) GetStore(store, name string) string {
	if strings.HasPrefix(store, qemuimgdriver.Location) {
		return store + "/" + name + "/"
	} else {
		return qemuimgdriver.Location + store + "/" + name + "/"
	}
}

func (h *Server) NewImage(conf *InstanceConfSchema) (*qemuimgdriver.Image, error) {
	path := h.GetStore(conf.DataStore, conf.Name)
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}
	file := path + "disk0.qcow2"
	size := libstar.ToBytes(conf.DiskSize, conf.DiskUnit)
	img := qemuimgdriver.NewImage(file, size)
	if err := img.Create(); err != nil {
		return nil, err
	}
	return img, nil
}

func (h *Server) InstanceConf2XML(conf *InstanceConfSchema) (libvirtdriver.DomainXML, error) {
	dom := libvirtdriver.DomainXML{
		Type: "kvm",
		Name: conf.Name,
		Devices: libvirtdriver.DevicesXML{
			Disks:    make([]libvirtdriver.DiskXML, 2),
			Graphics: make([]libvirtdriver.GraphicsXML, 1),
		},
		OS: libvirtdriver.OSXML{
			Type: libvirtdriver.OSTypeXML{
				Arch:  conf.Arch,
				Value: "hvm",
			},
			Boot: make([]libvirtdriver.OSBootXML, 3),
		},
	}
	if dom.OS.Type.Arch == "" {
		dom.OS.Type.Arch = "x86_64"
	}
	// create new disk firstly.
	img, err := h.NewImage(conf)
	if err != nil {
		return dom, err
	}
	if conf.Boots == "" {
		conf.Boots = "hd,cdrom,network"
	}
	for i, v := range strings.Split(conf.Boots, ",") {
		if i < 3 {
			dom.OS.Boot[i] = libvirtdriver.OSBootXML{
				Dev: v,
			}
		}
	}
	dom.VCPUXml = libvirtdriver.VCPUXML{
		Placement: "static",
		Value:     conf.Cpu,
	}
	dom.Memory = libvirtdriver.MemXML{
		Value: conf.MemorySize,
		Type:  conf.MemoryUnit,
	}
	dom.CurMem = libvirtdriver.CurMemXML{
		Value: conf.MemorySize,
		Type:  conf.MemoryUnit,
	}
	dom.Devices.Graphics[0] = libvirtdriver.GraphicsXML{
		Type:   "vnc",
		Listen: "0.0.0.0",
		Port:   "-1",
	}
	if strings.HasPrefix(conf.IsoFile, "/dev") {
		dom.Devices.Disks[0] = libvirtdriver.DiskXML{
			Type:   "block",
			Device: "cdrom",
			Driver: libvirtdriver.DiskDriverXML{
				Name: "qemu",
				Type: "raw",
			},
			Source: libvirtdriver.DiskSourceXML{
				Device: conf.IsoFile,
			},
			Target: libvirtdriver.DiskTargetXML{
				Bus: "ide",
				Dev: "hda",
			},
		}
	} else {
		dom.Devices.Disks[0] = libvirtdriver.DiskXML{
			Type:   "file",
			Device: "cdrom",
			Driver: libvirtdriver.DiskDriverXML{
				Name: "qemu",
				Type: "raw",
			},
			Source: libvirtdriver.DiskSourceXML{
				File: conf.IsoFile,
			},
			Target: libvirtdriver.DiskTargetXML{
				Bus: "ide",
				Dev: "hda",
			},
		}
	}
	dom.Devices.Disks[1] = libvirtdriver.DiskXML{
		Type:   "file",
		Device: "disk",
		Driver: libvirtdriver.DiskDriverXML{
			Name: "qemu",
			Type: img.Format,
		},
		Source: libvirtdriver.DiskSourceXML{
			File: img.Path,
		},
		Target: libvirtdriver.DiskTargetXML{
			Bus: "virtio",
			Dev: "vda",
		},
	}
	return dom, nil
}

func (h *Server) AddInstance(w http.ResponseWriter, r *http.Request) {
	hyper, err := libvirtdriver.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	xmlObj, err := h.InstanceConf2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	xmlData := xmlObj.Encode()
	if xmlData == "" {
		http.Error(w, "DomainXML.Encode has error.", http.StatusInternalServerError)
		return
	}
	file := h.GetStore(conf.DataStore, conf.Name) + "/define.xml"
	libstar.XML.MarshalSave(xmlObj, file, true)

	if dom, err := hyper.DomainDefineXML(xmlData); err == nil {
		if err := dom.Create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		domXML := libvirtdriver.NewDomainXMLFromDom(dom, true)
		if domXML != nil {
			h.ResponseJson(w, domXML)
		} else {
			h.ResponseJson(w, xmlObj)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Server) ModInstance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := h.GetArg(r, "id")

	hyper, err := libvirtdriver.GetHyper()
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
		xmlData, err := dom.GetXMLDesc(false)
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
		xmlData, err := dom.GetXMLDesc(false)
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
	case "destroy":
		if err := dom.Destroy(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "remove":
		if err := dom.Undefine(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if name, err := dom.GetName(); err == nil {
				path := h.GetStore("datastore/01", name)
				os.RemoveAll(path)
			}
		}
		return
	}
	h.ResponseMsg(w, 0, "")
}

func (h *Server) GetDataStore(w http.ResponseWriter, r *http.Request) {
	h.ResponseJson(w, qemuimgdriver.DATASTOR.List())
}

func (h *Server) GetISO(w http.ResponseWriter, r *http.Request) {
	store := h.GetQueryOne(r, "datastore")
	if store == "" {
		store = "datastore/01"
	}
	h.ResponseJson(w, qemuimgdriver.ISO.ListFiles(store))
}

func (h *Server) GetBridge(w http.ResponseWriter, r *http.Request) {
	h.ResponseJson(w, nil)
}
