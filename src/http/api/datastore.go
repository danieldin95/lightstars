package api

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/storage"
	"github.com/danieldin95/lightstar/src/storage/libvirts"
	"github.com/gorilla/mux"
	"github.com/libvirt/libvirt-go"

	//"github.com/libvirt/libvirt-go"
	"net/http"
	"sort"
)

type DataStore struct {
}

func DataStore2XML(conf schema.DataStore) libvirts.Pool {
	name := storage.PATH.GetStoreID(conf.Name)
	path := storage.PATH.Unix(conf.Name)

	xmlObj := &libvirts.PoolXML{
		Type: conf.Type,
		Name: name,
		Target: libvirts.TargetXML{
			Path: path,
		},
	}
	if conf.Type == "netfs" && conf.NFS != nil {
		xmlObj.Source = libvirts.SourceXML{
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
		XML:  libstar.XML.Encode(xmlObj),
	}
}

func (store DataStore) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore", store.Get).Methods("GET")
	router.HandleFunc("/api/datastore", store.Post).Methods("POST")
	router.HandleFunc("/api/datastore/{id}", store.Get).Methods("GET")
	router.HandleFunc("/api/datastore/{id}", store.Delete).Methods("DELETE")
}

func (store DataStore) Get(w http.ResponseWriter, r *http.Request) {
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

	pool, err := libvirts.LookupPoolByUUID(uuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pool.Free()
	format := GetQueryOne(r, "format")
	if format == "xml" {
		xmlDesc, err := pool.GetXMLDesc(libvirt.STORAGE_XML_INACTIVE)
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

func (store DataStore) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	if err := RemovePool(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}
