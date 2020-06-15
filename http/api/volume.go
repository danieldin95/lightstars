package api

import (
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type Volume struct {
}

func (v Volume) Router(router *mux.Router) {
	router.HandleFunc("/api/volume/{id}", v.GET).Methods("GET")
}

func (v Volume) Get(name string, data schema.VolumeInfos) error {

	return nil
}

func (v Volume) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	data := schema.VolumeInfos{
		Items: make([]schema.VolumeInfo, 0, 32),
	}
	infos, err := (&libvirts.Pool{Name: uuid}).List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, info := range infos {
		i := &schema.VolumeInfo{
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

