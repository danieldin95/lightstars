package http

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type UI struct {
}

func (ui UI) Router(router *mux.Router) {
	router.HandleFunc("/", ui.Index)
	router.HandleFunc("/ui", ui.Home)
	router.HandleFunc("/ui/", ui.Home)
	router.HandleFunc("/ui/index", ui.Home)
	router.HandleFunc("/ui/console", ui.Console)
	router.HandleFunc("/ui/instance/{id}", ui.Instance)
}

func (ui UI) Instance(w http.ResponseWriter, r *http.Request) {
	uuid, _ := api.GetArg(r, "id")

	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()
	instance := schema.NewInstance(*dom)
	file := api.GetFile("ui/instance.html")
	if err := api.ParseFiles(w, file, instance); err != nil {
		libstar.Error("UI.Instance %s", err)
	}
}

func (ui UI) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusTemporaryRedirect)
}

func (ui UI) Home(w http.ResponseWriter, r *http.Request) {
	index := schema.Index{
		Instances:  make([]schema.Instance, 0, 32),
		DataStores: make([]schema.DataStore, 0, 32),
		Networks:   make([]schema.Network, 0, 32),
	}

	index.Version = schema.NewVersion()
	index.Hyper = schema.NewHyper()
	if domains, err := libvirtc.ListDomains(); err == nil {
		for _, d := range domains {
			instance := schema.NewInstance(d)
			index.Instances = append(index.Instances, instance)
			d.Free()
		}
	}
	if pools, err := libvirts.ListPools(); err == nil {
		for _, p := range pools {
			store := schema.NewDataStore(p)
			index.DataStores = append(index.DataStores, store)
			p.Free()
		}
	}
	if nets, err := libvirtn.ListNetworks(); err == nil {
		for _, net := range nets {
			n := schema.NewNetwork(net)
			index.Networks = append(index.Networks, n)
			net.Free()
		}
	}
	file := api.GetFile("ui/index.html")
	if err := api.ParseFiles(w, file, index); err != nil {
		libstar.Error("UI.Home %s", err)
	}
}

func (ui UI) Console(w http.ResponseWriter, r *http.Request) {
	uuid := api.GetQueryOne(r, "instance")
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
	file := api.GetFile("ui/console.html")
	if err := api.ParseFiles(w, file, instance); err != nil {
		libstar.Error("UI.Console %s", err)
	}
}

func (ui UI) Hi(w http.ResponseWriter, r *http.Request) {
	id, _ := api.GetArg(r, "id")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	libstar.Info("UI.Hi id: %s, body: %s", id, body)
	api.ResponseJson(w, nil)
}
