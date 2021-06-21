package api

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/storage/libvirts"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

type Upload struct {
}

func (up Upload) Router(router *mux.Router) {
	router.HandleFunc("/api/upload/{id}", up.Post).Methods("POST")
	router.HandleFunc("/api/upload/{id}/volume/{name}", up.Get).Methods("GET")
}

func (up Upload) Get(w http.ResponseWriter, r *http.Request) {
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

func (up Upload) Post(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	pol, err := libvirts.LookupPoolByUUIDOrName(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer pol.Free()
	desc, err := pol.GetXMLDesc(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pool := &libvirts.PoolXML{}
	if err := libstar.XML.Decode(pool, desc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := pool.Target.Path
	// no more than 50MiB of memory
	_ = r.ParseMultipartForm(50 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		libstar.Error("Upload.Post %v", err)
		return
	}
	defer file.Close()
	modes := os.O_RDWR | os.O_CREATE | os.O_EXCL
	tempFile, err := os.OpenFile(path+"/"+handler.Filename, modes, 0660)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		libstar.Error("Upload.Post: %v", err)
		return
	}
	defer tempFile.Close()
	_, _ = io.Copy(tempFile, file)

	libstar.Info("Upload.Post saved to %s", path)
	ResponseMsg(w, 0, handler.Filename+" success")
}
