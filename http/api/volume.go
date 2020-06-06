package api

import (
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
)

type Volume struct {
}

func (v Volume) Router(router *mux.Router) {
	router.HandleFunc("/api/volume/{id}", v.GET).Methods("GET")
}

func (v Volume) Get(name string, data schema.VolumeInfos) error {
	infos, err := (&libvirts.Pool{Name: name}).List()
	if err != nil {
		return err
	}
	for name, info := range infos {
		data[name] = schema.VolumeInfo{
			Pool:       info.Pool,
			Name:       info.Name,
			Type:       info.Type,
			Allocation: info.Allocation,
			Capacity:   info.Capacity,
		}
	}
	return nil
}

func (v Volume) GET(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	data := make(schema.VolumeInfos, 128)
	if err := v.Get(uuid, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
