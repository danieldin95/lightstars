package api

import (
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/storage/libvirts"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"sort"
)

type Volume struct {
}

func (v Volume) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore/{id}/volume", v.GET).Methods("GET")
	router.HandleFunc("/api/datastore/{id}/volume/{name}", v.DOWNLOAD).Methods("GET")
	router.HandleFunc("/api/datastore/{id}/volume/{name}", v.DELETE).Methods("DELETE")
}

func (v Volume) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	data := schema.Volumes{
		Items: make([]schema.Volume, 0, 32),
	}
	pool := &libvirts.Pool{Name: uuid}
	infos, err := pool.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, info := range infos {
		i := &schema.Volume{
			Pool:       info.Pool,
			Name:       info.Name,
			Type:       info.Type,
			Allocation: info.Allocation,
			Capacity:   info.Capacity,
		}
		data.Items = append(data.Items, *i)
	}
	sort.Slice(data.Items, func(i, j int) bool {
		return data.Items[i].Name < data.Items[j].Name
	})
	data.Metadata.Size = len(data.Items)
	data.Metadata.Total = len(data.Items)
	ResponseJson(w, data)
}

func (v Volume) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (v Volume) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (v Volume) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	pool, err := libvirts.LookupPoolByUUID(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pool.Free()
	name, _ := GetArg(r, "name")
	vol, err := pool.LookupStorageVolByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer vol.Free()
	if err := vol.Delete(0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "")
}

func (v Volume) DOWNLOAD(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	pool, err := libvirts.LookupPoolByUUIDOrName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pool.Free()
	name, _ := GetArg(r, "name")
	vol, err := pool.LookupStorageVolByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	path, err := vol.GetPath()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = io.Copy(w, file)
}
