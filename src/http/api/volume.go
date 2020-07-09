package api

import (
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Volume struct {
}

func (v Volume) Router(router *mux.Router) {
	router.HandleFunc("/api/datastore/{id}/volume", v.GET).Methods("GET")
	router.HandleFunc("/api/datastore/{id}/volume/{name}", v.DOWNLOAD).Methods("GET")
}

func (v Volume) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	data := schema.Volumes{
		Items: make([]schema.Volume, 0, 32),
	}
	infos, err := (&libvirts.Pool{Name: uuid}).List()
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
	ResponseMsg(w, 0, "")
}

func (v Volume) DOWNLOAD(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
