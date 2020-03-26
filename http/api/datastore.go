package api

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type DataStore struct {
}

func DataStore2XML(conf schema.DataStore) libvirts.Pool {
	name := storage.PATH.GetStoreID(conf.Name)
	path := storage.PATH.Unix(conf.Name)

	polXml := libvirts.PoolXML{
		Type: conf.Type,
		Name: name,
		Target: libvirts.TargetXML{
			Path: path,
		},
	}
	if conf.Type == "netfs" && conf.NFS != nil {
		polXml.Source = libvirts.SourceXML{
			Host: libvirts.HostXML{
				Name: conf.NFS.Host,
			},
			Dir: libvirts.DirXML{
				Path: conf.NFS.Path,
			},
			Format: libvirts.FormatXML{
				Type: "nfs",
			},
		}
	}
	return libvirts.Pool{
		Type: conf.Type,
		Name: name,
		Path: path,
		XML:  polXml.Encode(),
	}
}

func (store DataStore) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore", store.GET).Methods("GET")
	router.HandleFunc("/api/datastore", store.POST).Methods("POST")
	router.HandleFunc("/api/datastore/{id}", store.DELETE).Methods("DELETE")
}

func (store DataStore) GET(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	if !ok {
		// list all instances.
		list := schema.ListDataStore{
			Items: make([]schema.DataStore, 0, 32),
		}
		if pools, err := libvirts.ListPools(); err == nil {
			for _, p := range pools {
				store := storage.NewDataStore(p)
				list.Items = append(list.Items, store)
				p.Free()
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
	//TODO
	ResponseJson(w, uuid)
}

func (store DataStore) POST(w http.ResponseWriter, r *http.Request) {
	data := schema.DataStore{}
	if err := GetData(r, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pol := DataStore2XML(data)
	libstar.Debug("DataStore.POST %s", pol.XML)
	if err := pol.Create(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (store DataStore) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, nil)
}

func (store DataStore) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	if err := libvirts.RemovePool(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}
