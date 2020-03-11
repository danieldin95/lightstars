package http

import (
	"encoding/base64"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/http/service"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/gorilla/mux"
	"net/http"
)

type Login struct {
}

func (l Login) Router(router *mux.Router) {
	router.HandleFunc("/login", l.Login)
}

func (l Login) Login(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Error string
	}{}

	if r.Method == "POST" {
		name := r.FormValue("name")
		pass := r.FormValue("password")
		u, _ := service.USERS.Get(name)
		if u.Password != pass {
			data.Error = "Invalid username or password."
		} else {
			basic := name + ":" + pass
			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: base64.StdEncoding.EncodeToString([]byte(basic)),
				Path:  "/",
			})
			http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
		}
	}
	file := api.GetFile("ui/login.html")
	if err := api.ParseFiles(w, file, data); err != nil {
		libstar.Error("Login.Instance %s", err)
	}
}
