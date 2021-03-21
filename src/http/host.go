package http

import (
	"github.com/danieldin95/lightstar/src/http/api"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/service"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strings"
)

type Host struct {
}

func (h Host) Router(router *mux.Router) {
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("GET")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("POST")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("PUT")
	router.PathPrefix("/host/{name}/api/").HandlerFunc(h.Handle).Methods("DELETE")
	router.PathPrefix("/host/{name}/ext/").HandlerFunc(h.Handle).Methods("GET")
}

func (h Host) filterInstance(r *http.Response, w http.ResponseWriter, user *schema.User) {
	all := schema.ListInstance{
		Items: make([]schema.Instance, 0, 32),
	}
	if err := libstar.GetJSON(r.Body, &all); err != nil {
		libstar.Warn("Host.Filter %s", r.Request.URL)
		return
	}
	list := schema.ListInstance{
		Items: make([]schema.Instance, 0, 32),
	}
	obj := api.Instance{}
	for _, item := range all.Items {
		if obj.HasPermission(user, item.Name) {
			list.Items = append(list.Items, item)
		}
	}
	sort.SliceStable(list.Items, func(i, j int) bool {
		return list.Items[i].Name < list.Items[j].Name
	})
	list.Metadata.Size = len(list.Items)
	list.Metadata.Total = len(list.Items)
	api.ResponseJson(w, list)
}

func (h Host) filterPort(r *http.Response, w http.ResponseWriter, user *schema.User) {
	all := schema.ListInterface{
		Items: make([]schema.Interface, 0, 32),
	}
	if err := libstar.GetJSON(r.Body, &all); err != nil {
		libstar.Warn("Host.Filter %s", r.Request.URL)
		return
	}
	list := schema.ListInterface{
		Items: make([]schema.Interface, 0, 32),
	}
	obj := api.Instance{}
	for _, item := range all.Items {
		if obj.HasPermission(user, item.Domain.Name) {
			list.Items = append(list.Items, item)
		}
	}
	list.Metadata.Size = len(list.Items)
	list.Metadata.Total = len(list.Items)
	api.ResponseJson(w, list)
}

func (h Host) Filter(r *http.Response, w http.ResponseWriter, data interface{}) bool {
	req := r.Request
	if data == nil || req == nil || req.Method != "GET" {
		return false
	}
	user := data.(*schema.User)

	if req.Method == "GET" {
		libstar.Debug("Host.Filter %s %s %s", user.Name, req.Method, req.URL.Path)
	} else {
		libstar.Info("Host.Filter %s %s %s", user.Name, req.Method, req.URL.Path)
	}

	u := strings.Split(req.URL.Path, "/")
	if len(u) == 3 && u[1] == "api" && u[2] == "instance" {
		h.filterInstance(r, w, user)
	}
	if len(u) == 5 && u[1] == "api" && u[2] == "network" && u[4] == "interface" {
		h.filterPort(r, w, user)
	} else {
		return false
	}
	return true
}

func (h Host) Handle(w http.ResponseWriter, r *http.Request) {
	user, _ := api.GetUser(r)
	name, _ := api.GetArg(r, "name")
	node := service.SERVICE.Zone.Get(name)
	if node == nil {
		http.Error(w, "host not found", http.StatusNotFound)
		return
	}
	libstar.Debug("Host.Handle %s", node)
	r.Header.Del("cookie")
	pri := &libstar.ProxyUrl{
		Proxy: libstar.Proxy{
			Prefix: "/host/" + name,
			Server: node.Url,
			Auth: libstar.Auth{
				Type:     "basic",
				Username: node.Username,
				Password: node.Password,
			},
		},
		Filter: h.Filter,
		Data:   &user,
	}
	pri.Initialize()
	pri.Handler(w, r)
}
