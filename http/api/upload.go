package api

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

type Upload struct {
}

func (up Upload) Router(router *mux.Router) {
	router.HandleFunc("/api/upload/{id}", up.Upload).Methods("POST")
}

func (up Upload) Upload(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")

	pol, err := libvirts.LookupPoolByUUID(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer pol.Free()

	name, err := pol.GetName()
	path := storage.PATH.Unix(name)

	// no more than 50MiB of memory
	_ = r.ParseMultipartForm(50 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		libstar.Error("Upload.Upload %v", err)
		return
	}
	defer file.Close()
	modes := os.O_RDWR | os.O_CREATE | os.O_EXCL
	tempFile, err := os.OpenFile(path+"/"+handler.Filename, modes, 0660)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		libstar.Error("Upload.Upload: %v", err)
		return
	}
	defer tempFile.Close()
	_, _ = io.Copy(tempFile, file)

	libstar.Info("Upload.Upload saved to %s", path)
	ResponseMsg(w, 0, "success")
}
