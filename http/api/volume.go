package api

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"

	//"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"path"
)


type Volume struct {
	Name string
	State string
	Autostart	bool
	Persistent	bool
	Capacity	string
	Allocation	string
	Available	string
}

func (v Volume) Router(router *mux.Router) {
	router.PathPrefix("/api/datastore/{id}/vol").HandlerFunc(v.GET).Methods("GET")
}


func (v Volume) GET(w http.ResponseWriter, r *http.Request)  {
	var store string
	uuid, ok := GetArg(r, "id")
	if !ok {
		http.Error(w, "store uuid not found", http.StatusNotFound)
		return
	}

	pool, err := libvirts.LookupPoolByUUID(uuid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pool.Free()
	libstar.Info("store is  %s", storage.NewDataStore(*pool).Name)
	store = storage.NewDataStore(*pool).Name
	p := storage.PATH.Unix(store)

	dir := GetQueryOne(r, "dir")
	if dir != "" {
		p = path.Join(p, dir)
	}

	libstar.Info("file is  %s %s", p, store)
	ResponseJson(w, storage.FILE.List(p))
}

func (v Volume) POST(w http.ResponseWriter, r *http.Request)  {
	ResponseMsg(w, 0, "")
}

func (v Volume) PUT(w http.ResponseWriter, r *http.Request)  {
	ResponseMsg(w, 0, "")
}

func (v Volume) DELETE(w http.ResponseWriter, r *http.Request)  {
	ResponseMsg(w, 0, "")
}

