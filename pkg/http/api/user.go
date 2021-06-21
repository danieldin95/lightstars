package api

import (
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

type User struct {
}

type Password struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func (u User) Router(router *mux.Router) {
	router.HandleFunc("/api/user", u.Get).Methods("GET")
	router.HandleFunc("/api/user", u.Post).Methods("POST")
	router.HandleFunc("/api/user/{name}/password", u.Put).Methods("PUT")
	router.HandleFunc("/api/user/{name}", u.Delete).Methods("DELETE")
}

func (u User) Get(writer http.ResponseWriter, request *http.Request) {
	users := schema.ListUser{
		Items: make([]schema.User, 0, 32),
	}
	for user := range service.SERVICE.Users.List() {
		if user == nil {
			break
		}
		users.Items = append(users.Items, *user)
	}
	sort.SliceStable(users.Items, func(i, j int) bool {
		return users.Items[i].Name < users.Items[j].Name
	})
	ResponseJson(writer, users)
}

func (u User) Post(writer http.ResponseWriter, request *http.Request) {
	data := &schema.User{}
	if err := GetData(request, data); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err := service.SERVICE.Users.Add(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	_ = service.SERVICE.Users.Save()
	ResponseMsg(writer, 0, "success")
}

func (u User) Put(writer http.ResponseWriter, request *http.Request) {
	data := &Password{}
	if err := GetData(request, data); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	name, _ := GetArg(request, "name")
	_, ok := service.SERVICE.Users.SetPass(name, data.Old, data.New)
	if !ok {
		http.Error(writer, "invalid password", http.StatusBadRequest)
		return
	}
	_ = service.SERVICE.Users.Save()
	ResponseMsg(writer, 0, "success")
}

func (u User) Delete(writer http.ResponseWriter, request *http.Request) {
	name, _ := GetArg(request, "name")
	err := service.SERVICE.Users.Del(name)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	ResponseMsg(writer, 0, "success")
}
