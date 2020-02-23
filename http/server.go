package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
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
	if h.router != nil {
		return h.router
	}
	h.router = mux.NewRouter()
	h.router.NotFoundHandler = http.HandlerFunc(h.Handle404)
	h.router.Use(h.Middleware)
	return h.router
}

func (h *Server) LoadRouter() {
	router := h.Router()
	// static files
	staticFile := http.StripPrefix("/static/", http.FileServer(http.Dir(h.pubDir)))
	router.PathPrefix("/static/").Handler(staticFile)
	// proxy websocket
	router.Handle("/websockify", websocket.Handler(h.HandleSockify))
	// custom router
	router.HandleFunc("/", h.HandleIndex)
	router.HandleFunc("/ui", h.HandleUi)
	router.HandleFunc("/ui/", h.HandleUi)
	router.HandleFunc("/ui/index", h.HandleUi)
	router.HandleFunc("/ui/console", h.HandleConsole)
	router.HandleFunc("/ui/instance/{id}", h.HandleInstance)
	// api router
	router.HandleFunc("/api/instance", h.AddInstance).Methods("POST")
	router.HandleFunc("/api/instance/{id}", h.GetInstance).Methods("GET")
	router.HandleFunc("/api/instance/{id}", h.ModInstance).Methods("PUT")
	router.HandleFunc("/api/instance/{id}", h.DelInstance).Methods("DELETE")
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

func (h *Server) ResponseXML(w http.ResponseWriter, v string) {
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
	return fmt.Sprintf("%s/%s", h.pubDir, name)
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

func (h *Server) HandleFile(w http.ResponseWriter, r *http.Request) {
	realpath := h.GetFile(r.URL.Path)
	if _, err := os.Stat(realpath); !os.IsExist(err) {
		realpath = realpath + ".html"
	}
	contents, err := ioutil.ReadFile(realpath)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	fmt.Fprintf(w, "%s\n", contents)
}

func (h *Server) Handle404(w http.ResponseWriter, r *http.Request) {
	file := h.GetFile("404.html")
	if err := h.ParseFiles(w, file, nil); err != nil {
		libstar.Error("Server.Handle404 %s", err)
	}
}

func (h *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusTemporaryRedirect)
}

func (h *Server) HandleUi(w http.ResponseWriter, r *http.Request) {
	index := schema.Index{
		Instances: make([]schema.Instance, 0, 32),
	}

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		libstar.Error("Server.HandleIndex %s", err)
		return
	}
	index.Version = schema.NewVersion()
	index.Hyper = schema.NewHyper()
	if domains, err := hyper.ListAllDomains(); err == nil {
		for _, dom := range domains {
			instance := schema.NewInstance(dom)
			index.Instances = append(index.Instances, instance)
			dom.Free()
		}
	}
	file := h.GetFile("ui/index.html")
	if err := h.ParseFiles(w, file, index); err != nil {
		libstar.Error("Server.HandleIndex %s", err)
	}
}

func (h *Server) HandleConsole(w http.ResponseWriter, r *http.Request) {
	uuid := h.GetQueryOne(r, "instance")
	if uuid == "" {
		http.Error(w, "Not found instance", http.StatusNotFound)
		return
	}
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	instance := schema.NewInstance(*dom)
	file := h.GetFile("ui/console.html")
	if err := h.ParseFiles(w, file, instance); err != nil {
		libstar.Error("Server.HandleInstance %s", err)
	}
}

func (h *Server) HandleInstance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := h.GetArg(r, "id")

	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	instance := schema.NewInstance(*dom)
	file := h.GetFile("ui/instance.html")
	if err := h.ParseFiles(w, file, instance); err != nil {
		libstar.Error("Server.HandleInstance %s", err)
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

	libstar.Debug("Server.GetTarget %s", id)
	hyper, err := libvirtc.GetHyper()
	if err != nil {
		libstar.Error("Server.HandleIndex %s", err)
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

func (h *Server) HandleSockify(ws *websocket.Conn) {
	defer ws.Close()
	ws.PayloadType = websocket.BinaryFrame

	target := h.GetTarget(ws.Request())
	if target == "" {
		libstar.Error("Server.HandleSockify target not found.")
		return
	}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		libstar.Error("Server.HandleSockify dial %s", err)
		return
	}
	defer conn.Close()
	libstar.Info("Server.HandleSockify request from %s", ws.RemoteAddr())
	libstar.Info("Server.HandleSockify connect to %s", conn.RemoteAddr())

	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		if _, err := io.Copy(conn, ws); err != nil {
			libstar.Error("Server.HandleSockify copy from ws %s", err)
		}
	}()
	go func() {
		defer wait.Done()
		if _, err := io.Copy(ws, conn); err != nil {
			libstar.Error("Server.HandleSockify copy from target %s", err)
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

	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	format := h.GetQueryOne(r, "format")
	if format == "xml" {
		xmlDesc, err := dom.GetXMLDesc(false)
		if err == nil {
			h.ResponseXML(w, xmlDesc)
		} else {
			h.ResponseXML(w, "<error>"+err.Error()+"</error>")
		}
	} else if format == "schema" {
		h.ResponseJson(w, schema.NewInstance(*dom))
	} else {
		h.ResponseJson(w, libvirtc.NewDomainXMLFromDom(dom, true))
	}
}

func (h *Server) GetPath(store, name string) string {
	return storage.PATH.Unix(store) + "/" + name + "/"
}

func (h *Server) NewVolumeAndPool(conf *schema.InstanceConf) (*libvirts.VolumeXML, error) {
	path := h.GetPath(conf.DataStore, conf.Name)
	pol, err := libvirts.CreatePool(libvirts.ToDomainPool(conf.Name), path)
	if err != nil {
		return nil, err
	}
	size := libstar.ToBytes(conf.DiskSize, conf.DiskUnit)
	vol, err := libvirts.CreateVolume(pol.Name, "disk0.qcow2", size)
	if err != nil {
		return nil, err
	}
	return vol.GetXMLObj()
}

func (h *Server) DelVolumeAndPool(name string) error {
	err := libvirts.RemovePool(libvirts.ToDomainPool(name))
	if err != nil {
		return err
	}
	return nil
}

func (h *Server) InstanceConf2XML(conf *schema.InstanceConf) (libvirtc.DomainXML, error) {
	dom := libvirtc.DomainXML{
		Type: "kvm",
		Name: conf.Name,
		Devices: libvirtc.DevicesXML{
			Disks:      make([]libvirtc.DiskXML, 2),
			Graphics:   make([]libvirtc.GraphicsXML, 1),
			Interfaces: make([]libvirtc.InterfaceXML, 1),
		},
		OS: libvirtc.OSXML{
			Type: libvirtc.OSTypeXML{
				Arch:  conf.Arch,
				Value: "hvm",
			},
			Boot: make([]libvirtc.OSBootXML, 3),
		},
	}
	if dom.OS.Type.Arch == "" {
		dom.OS.Type.Arch = "x86_64"
	}
	// create new disk firstly.
	vol, err := h.NewVolumeAndPool(conf)
	if err != nil {
		return dom, err
	}
	if conf.Boots == "" {
		conf.Boots = "hd,cdrom,network"
	}
	for i, v := range strings.Split(conf.Boots, ",") {
		if i < 3 {
			dom.OS.Boot[i] = libvirtc.OSBootXML{
				Dev: v,
			}
		}
	}
	dom.VCPUXml = libvirtc.VCPUXML{
		Placement: "static",
		Value:     conf.Cpu,
	}
	dom.Memory = libvirtc.MemXML{
		Value: conf.MemorySize,
		Type:  conf.MemoryUnit,
	}
	dom.CurMem = libvirtc.CurMemXML{
		Value: conf.MemorySize,
		Type:  conf.MemoryUnit,
	}
	dom.Devices.Graphics[0] = libvirtc.GraphicsXML{
		Type:   "vnc",
		Listen: "0.0.0.0",
		Port:   "-1",
	}
	if strings.HasPrefix(conf.IsoFile, "/dev") {
		dom.Devices.Disks[0] = libvirtc.DiskXML{
			Type:   "block",
			Device: "cdrom",
			Driver: libvirtc.DiskDriverXML{
				Name: "qemu",
				Type: "raw",
			},
			Source: libvirtc.DiskSourceXML{
				Device: conf.IsoFile,
			},
			Target: libvirtc.DiskTargetXML{
				Bus: "ide",
				Dev: "hda",
			},
		}
	} else {
		dom.Devices.Disks[0] = libvirtc.DiskXML{
			Type:   "file",
			Device: "cdrom",
			Driver: libvirtc.DiskDriverXML{
				Name: "qemu",
				Type: "raw",
			},
			Source: libvirtc.DiskSourceXML{
				File: storage.PATH.Unix(conf.IsoFile),
			},
			Target: libvirtc.DiskTargetXML{
				Bus: "ide",
				Dev: "hda",
			},
		}
	}
	dom.Devices.Disks[1] = libvirtc.DiskXML{
		Type:   "file",
		Device: "disk",
		Driver: libvirtc.DiskDriverXML{
			Name: "qemu",
			Type: vol.Target.Format.Type,
		},
		Source: libvirtc.DiskSourceXML{
			File: vol.Target.Path,
		},
		Target: libvirtc.DiskTargetXML{
			Bus: "virtio",
			Dev: "vda",
		},
	}
	dom.Devices.Interfaces[0] = libvirtc.InterfaceXML{
		Type: "bridge",
		Source: libvirtc.InterfaceSourceXML{
			Bridge: conf.Interface,
		},
		Model: libvirtc.InterfaceModelXML{
			Type: "virtio",
		},
	}
	return dom, nil
}

func (h *Server) AddInstance(w http.ResponseWriter, r *http.Request) {
	hyper, err := libvirtc.GetHyper()
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
	conf := &schema.InstanceConf{}
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
	file := storage.PATH.RootXML() + conf.Name + ".xml"
	libstar.XML.MarshalSave(xmlObj, file, true)

	dom, err := hyper.DomainDefineXML(xmlData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()
	if conf.Start == "true" {
		if err := dom.Create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	domXML := libvirtc.NewDomainXMLFromDom(dom, true)
	if domXML != nil {
		h.ResponseJson(w, domXML)
	} else {
		h.ResponseJson(w, xmlObj)
	}
}

func (h *Server) ModInstance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := h.GetArg(r, "id")

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conf := &schema.InstanceConf{}
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
		if domNew, err := hyper.DomainDefineXML(xmlData); err == nil {
			if err := dom.Create(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			domNew.Free()
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "shutdown":
		xmlData, err := dom.GetXMLDesc(false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := dom.Shutdown(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		domNew, err := hyper.DomainDefineXML(xmlData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		domNew.Free()
	case "suspend":
		if err := dom.Suspend(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "destroy":
		if err := dom.Destroy(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "resume":
		if err := dom.Resume(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "undefine":
		if err := dom.Undefine(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	h.ResponseMsg(w, 0, "success")
}

func (h *Server) DelInstance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := h.GetArg(r, "id")

	hyper, err := libvirtc.GetHyper()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dom, err := hyper.LookupDomainByUUIDName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()

	if ok, _ := dom.IsActive(); ok {
		http.Error(w, "not allowed with active instance", http.StatusInternalServerError)
		return
	}
	name, err := dom.GetName()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := dom.Undefine(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.DelVolumeAndPool(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.ResponseMsg(w, 0, "")
}

func (h *Server) GetDataStore(w http.ResponseWriter, r *http.Request) {
	h.ResponseJson(w, libvirts.DATASTOR.List())
}

func (h *Server) GetISO(w http.ResponseWriter, r *http.Request) {
	store := h.GetQueryOne(r, "datastore")
	if store == "" {
		store = "datastore@01"
	}
	path := storage.PATH.Unix(store)
	h.ResponseJson(w, libvirts.ISO.ListFiles(path))
}

func (h *Server) GetBridge(w http.ResponseWriter, r *http.Request) {
	h.ResponseJson(w, nil)
}
