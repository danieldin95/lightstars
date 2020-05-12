package api

import (
	"github.com/danieldin95/lightstar/service"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {

}

type Password struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (u User) Router(router *mux.Router) {
	router.HandleFunc("/api/user/password", u.PUT).Methods("POST")
}

func (u User) PUT(writer http.ResponseWriter, request *http.Request) {
	data  := &Password{}
	if err := GetData(request, data); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

        name, _, _ := GetAuth(request)
        
	if _, status := service.SERVICE.Users.SetPassWord(name, data.OldPassword, data.NewPassword); status != true {
		ResponseMsg(writer, 200, "password error")
                return
	} 
	service.SERVICE.Users.Save()
        ResponseMsg(writer, 0, "success")
}

