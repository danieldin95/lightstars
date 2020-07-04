package api

import (
	"github.com/danieldin95/lightstar/src/service"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
}

type Password struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func (u User) Router(router *mux.Router) {
	router.HandleFunc("/api/user/password", u.PUT).Methods("POST")
}

func (u User) PUT(writer http.ResponseWriter, request *http.Request) {
	data := &Password{}
	if err := GetData(request, data); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	name, _, _ := GetAuth(request)
	_, ok := service.SERVICE.Users.SetPass(name, data.Old, data.New)
	if !ok {
		http.Error(writer, "invalid password", http.StatusBadRequest)
		return
	}
	_ = service.SERVICE.Users.Save()
	ResponseMsg(writer, 0, "success")
}
