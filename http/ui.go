package http

import (
	"encoding/base64"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"github.com/danieldin95/lightstar/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type UI struct {
}

func (ui UI) Router(router *mux.Router) {
	router.HandleFunc("/", ui.Index)
	router.HandleFunc("/ui", ui.Home)
	router.HandleFunc("/ui/", ui.Index)
	router.HandleFunc("/ui/index", ui.Index)
	router.HandleFunc("/ui/console", ui.Console)
	router.HandleFunc("/ui/spice", ui.Spice)
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
	index.Hyper = schema.NewHyper()
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

func (ui UI) Spice(w http.ResponseWriter, r *http.Request) {
	uuid := api.GetQueryOne(r, "id")
	if uuid == "" {
		http.Error(w, "Not found instance", http.StatusNotFound)
		return
	}
	file := api.GetFile("ui/spice.html")
	if err := api.ParseFiles(w, file, nil); err != nil {
		libstar.Error("UI.Spice %s", err)
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
		u, ok := service.SERVICE.Users.Get(name)
		if !ok || u.Password != pass {
			data.Error = "Invalid username or password."
		} else {
			basic := name + ":" + pass
			cookie := http.Cookie{
				Name:  "token",
				Value: base64.StdEncoding.EncodeToString([]byte(basic)),
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
			return
		}
	}
	file := api.GetFile("ui/login.html")
	if err := api.ParseFiles(w, file, data); err != nil {
		libstar.Error("Login.Instance %s", err)
	}
}
