package api

import (
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type DataStore struct {
}

func (store DataStore) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore", store.GET).Methods("GET")
	router.HandleFunc("/api/datastore", store.POST).Methods("POST")
	router.HandleFunc("/api/datastore/{id}", store.DELETE).Methods("DELETE")
}

func (store DataStore) GET(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	format := GetQueryOne(r, "format")
	if !ok {
		// list all instances.
		if format == "schema" {
			list := schema.List{
				Items:    make([]interface{}, 0, 32),
				Metadata: schema.MetaData{},
			}
			if pools, err := libvirts.ListPools(); err == nil {
				for _, p := range pools {
					store := schema.NewDataStore(p)
					list.Items = append(list.Items, store)
					p.Free()
				}
				sort.SliceStable(list.Items, func(i, j int) bool {
					return list.Items[i].(schema.DataStore).Name < list.Items[j].(schema.DataStore).Name
				})
				list.Metadata.Size = len(list.Items)
				list.Metadata.Total = len(list.Items)
			}
			ResponseJson(w, list)
		} else {
			ResponseJson(w, libvirts.DATASTOR.List())
		}
		return
	}

	//TODO
	ResponseJson(w, uuid)
}

func (store DataStore) POST(w http.ResponseWriter, r *http.Request) {
	data := &schema.DataStore{}
	if err := GetData(r, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if storage.PATH.IsDataStore(data.Name) {
		name := storage.PATH.GetStoreID(data.Name)
		path := storage.PATH.Unix(data.Name)
		if _, err := libvirts.CreatePool(name, path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
