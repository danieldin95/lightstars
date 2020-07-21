package http

import (
	"encoding/base64"
	"github.com/danieldin95/lightstar/src/compute"
	"github.com/danieldin95/lightstar/src/http/api"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/danieldin95/lightstar/src/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

type UI struct {
}

func (ui UI) Router(router *mux.Router) {
	router.HandleFunc("/", ui.Index)
	router.HandleFunc("/ui", ui.Home)
	router.HandleFunc("/ui/", ui.Index)
	router.HandleFunc("/ui/index", ui.Index)
	router.HandleFunc("/ui/console", ui.Console)
}

func (ui UI) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusTemporaryRedirect)
}

func (ui UI) Home(w http.ResponseWriter, r *http.Request) {
	index := schema.Index{}
	deft := service.SERVICE.Zone.Get("default")
	if deft != nil && deft.Url != "" {
		index.Default = "default"
	}
	index.User, _ = api.GetUser(r)
	libstar.Debug("UI.Home %v", index.User)
	index.Version = schema.NewVersion()
	index.Hyper = compute.NewHyper()
	file := api.GetFile("ui/index.html")
	if err := api.ParseFiles(w, file, index); err != nil {
		libstar.Error("UI.Home %s", err)
	}
}

func (ui UI) Console(w http.ResponseWriter, r *http.Request) {
	uuid := api.GetQueryOne(r, "id")
	if uuid == "" {
		http.Error(w, "Not found instance", http.StatusNotFound)
		return
	}
	file := api.GetFile("ui/console.html")
	if err := api.ParseFiles(w, file, nil); err != nil {
		libstar.Error("UI.Console %s", err)
	}
}

func (ui UI) Hi(w http.ResponseWriter, r *http.Request) {
	id, _ := api.GetArg(r, "id")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	libstar.Info("UI.Hi id: %s, body: %s", id, body)
	api.ResponseJson(w, nil)
}

type Login struct {
}

func (l Login) Router(router *mux.Router) {
	router.HandleFunc("/ui/login", l.Login)
}

func (l Login) Login(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Error string
	}{}
	if r.Method == "POST" {
		name := r.FormValue("name")
		pass := r.FormValue("password")
		next := r.FormValue("next")
		u, ok := service.SERVICE.Users.Get(name)
		if !ok || u.Password != pass {
			data.Error = "Invalid username or password."
		} else {
			basic := name + ":" + pass
			expired := time.Now().Add(time.Hour * 8)
			http.SetCookie(w, &http.Cookie{
				Name:    "token-id",
				Value:   base64.StdEncoding.EncodeToString([]byte(basic)),
				Path:    "/",
				Expires: expired,
			})
			uuid := libstar.GenToken(32)
			http.SetCookie(w, &http.Cookie{
				Name:    "session-id",
				Value:   uuid,
				Path:    "/",
				Expires: expired,
			})
			http.Redirect(w, r, "/ui#"+next, http.StatusMovedPermanently)
			return
		}
	}
	file := api.GetFile("ui/login.html")
	if err := api.ParseFiles(w, file, data); err != nil {
		libstar.Error("Login.Instance %s", err)
	}
}
