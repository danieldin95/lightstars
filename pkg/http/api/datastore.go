package api

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/storage"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type DataStore struct {
}

func DataStore2XML(conf schema.DataStore) storage.Pool {
	name := storage.PATH.GetStoreID(conf.Name)
	path := storage.PATH.Unix(conf.Name)

	xmlObj := &storage.PoolXML{
		Type: conf.Type,
		Name: name,
		Target: storage.TargetXML{
			Path: path,
		},
	}
	if conf.Type == "netfs" && conf.NFS != nil {
		xmlObj.Source = storage.SourceXML{
			Host: storage.HostXML{
				Name: conf.NFS.Host,
			},
			Dir: storage.DirXML{
				Path: conf.NFS.Path,
			},
			Format: storage.FormatXML{
				Type: "nfs",
			},
		}
	}
	return storage.Pool{
		Type: conf.Type,
		Name: name,
		Path: path,
		XML:  libstar.XML.Encode(xmlObj),
	}
}

func (store DataStore) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore", store.Get).Methods("GET")
	router.HandleFunc("/api/datastore", store.Post).Methods("POST")
	router.HandleFunc("/api/datastore/{id}", store.Get).Methods("GET")
	router.HandleFunc("/api/datastore/{id}/start", store.Start).Methods("PUT")
	router.HandleFunc("/api/datastore/{id}/destroy", store.Destroy).Methods("PUT")
	router.HandleFunc("/api/datastore/{id}/refresh", store.Refresh).Methods("PUT")
	router.HandleFunc("/api/datastore/{id}/autostart", store.Autostart).Methods("PUT")
	router.HandleFunc("/api/datastore/{id}/clean", store.Clean).Methods("PUT")
	router.HandleFunc("/api/datastore/{id}/remove", store.Remove).Methods("PUT")
	router.HandleFunc("/api/datastore/{id}", store.Delete).Methods("DELETE")
}

func (store DataStore) Get(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	if !ok {
		// list all instances.
		list := schema.ListDataStore{
			Items: make([]schema.DataStore, 0, 32),
		}
		if pools, err := storage.ListPools(); err == nil {
			for _, p := range pools {
				store := storage.NewDataStore(p)
				list.Items = append(list.Items, store)
				_ = p.Free()
			}
			sort.SliceStable(list.Items, func(i, j int) bool {
				return list.Items[i].Name < list.Items[j].Name
			})
			list.Metadata.Size = len(list.Items)
			list.Metadata.Total = len(list.Items)
		}
		ResponseJson(w, list)
		return
	}

	pool, err := storage.LookupPoolByUUID(uuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pool.Free()
	format := GetQueryOne(r, "format")
	if format == "xml" {
		xmlDesc, err := pool.GetXMLDesc(1)
		if err == nil {
			ResponseXML(w, xmlDesc)
		} else {
			ResponseXML(w, "<error>"+err.Error()+"</error>")
		}
	} else {
		ResponseJson(w, storage.NewDataStore(*pool))
	}
}

func (store DataStore) Post(w http.ResponseWriter, r *http.Request) {
	data := schema.DataStore{}
	if err := GetData(r, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pol := DataStore2XML(data)
	libstar.Debug("DataStore.Post %s", pol.XML)
	if err := pol.Create(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (store DataStore) Put(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (store DataStore) execute(w http.ResponseWriter, uuid, action string) {
	pol, err := storage.LookupPoolByUUIDOrName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pol.Free()
	switch action {
	case "start":
		if err := pol.Start(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "destroy":
		if err := pol.Destroy(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "refresh":
		if err := pol.Refresh(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "clean":
		if err := pol.Clean(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "remove":
		if err := pol.Remove(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "unsupported action", http.StatusBadRequest)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (store DataStore) Start(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	store.execute(w, uuid, "start")
}

func (store DataStore) Destroy(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	store.execute(w, uuid, "destroy")
}

func (store DataStore) Refresh(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	store.execute(w, uuid, "refresh")
}

func (store DataStore) Autostart(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	enable := GetQueryOne(r, "enable")
	on := !(enable == "false" || enable == "0" || enable == "no" || enable == "off")
	pol, err := storage.LookupPoolByUUIDOrName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pol.Free()
	if err := pol.SetAutostart(on); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (store DataStore) Clean(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	store.execute(w, uuid, "clean")
}

func (store DataStore) Remove(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	store.execute(w, uuid, "remove")
}

func (store DataStore) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	if err := RemovePool(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}
