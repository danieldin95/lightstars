package api

import (
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
)

type History struct {
}

func (his History) Router(router *mux.Router) {
	router.HandleFunc("/api/history", his.Get).Methods("GET")
}

func (his History) Get(w http.ResponseWriter, r *http.Request) {
	items := schema.ListHistory{
		Items: make([]schema.History, 0, 32),
	}
	user, _ := GetUser(r)
	name := user.Name
	if user.Type == "admin" {
		name = ""
	}
	for obj := range service.SERVICE.History.List(name) {
		if obj == nil {
			break
		}
		items.Items = append(items.Items, *obj)
	}
	items.Metadata.Size = len(items.Items)
	items.Metadata.Total = len(items.Items)
	ResponseJson(w, items)
}
