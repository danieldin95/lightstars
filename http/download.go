package http

import (
	"github.com/danieldin95/lightstar/storage"
	"github.com/gorilla/mux"
	"net/http"
)

type Download struct {
}

func (down Download) Router(router *mux.Router) {
	dir := http.Dir(storage.PATH.Root())
	files := http.StripPrefix("/download", http.FileServer(dir))
	router.PathPrefix("/download/").Handler(files)
}
