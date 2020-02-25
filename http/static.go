package http

import (
	"fmt"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

type Static struct {

}

func (s Static) Router(router *mux.Router) {
	dir := http.Dir(api.GetStatic())
	staticFile := http.StripPrefix("/static/", http.FileServer(dir))
	router.PathPrefix("/static/").Handler(staticFile)
}

func (s Static) HandleFile(w http.ResponseWriter, r *http.Request) {
	realpath := api.GetFile(r.URL.Path)
	if _, err := os.Stat(realpath); !os.IsExist(err) {
		realpath = realpath + ".html"
	}

	contents, err := ioutil.ReadFile(realpath)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	fmt.Fprintf(w, "%s\n", contents)
}
