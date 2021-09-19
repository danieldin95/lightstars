package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
	"github.com/danieldin95/lightstar/pkg/service"
	"github.com/danieldin95/lightstar/pkg/storage"
	"github.com/danieldin95/lightstar/pkg/storage/libvirts"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"text/template"
	"time"
)

func GetArg(r *http.Request, name string) (string, bool) {
	vars := mux.Vars(r)
	value, ok := vars[name]
	return value, ok
}

func GetData(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(body), v); err != nil {
		return err
	}
	return nil
}

func GetQueryOne(req *http.Request, name string) string {
	query := req.URL.Query()
	if values, ok := query[name]; ok {
		return values[0]
	}
	return ""
}

func ResponseJson(w http.ResponseWriter, v interface{}) {
	str, err := json.Marshal(v)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(str)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResponseXML(w http.ResponseWriter, v string) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = w.Write([]byte(v))
}

func ResponseMsg(w http.ResponseWriter, code int, message string) {
	if code == 0 && message == "" {
		message = "success"
	}
	ret := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	}
	ResponseJson(w, ret)
}

// store: datastore@01
// name: domain name
func GetPath(store, name string) string {
	return storage.PATH.Unix(store) + "/" + name + "/"
}

// name: domain name
// disk: disk name
// size: disk size using bytes
func NewVolume(name, disk string, size uint64) (*libvirts.VolumeXML, error) {
	vol, err := libvirts.CreateVolume(libvirts.ToDomainPool(name), disk, size)
	if err != nil {
		return nil, err
	}
	return vol.GetXMLObj()
}

// name: Domain name.
// store: like: datatore@01
func NewVolumeAndPool(store, name, disk string, size uint64) (*libvirts.VolumeXML, error) {
	path := GetPath(store, name)
	pol, err := libvirts.CreatePool(libvirts.ToDomainPool(name), path)
	if err != nil {
		return nil, err
	}
	vol, err := libvirts.CreateVolume(pol.Name, disk, size)
	if err != nil {
		return nil, err
	}
	return vol.GetXMLObj()
}

// name: Domain name.
// store: like: datatore@01
func NewBackingVolumeAndPool(store, name, disk, backingFle, backingFmt string) (*libvirts.VolumeXML, error) {
	path := GetPath(store, name)
	pol, err := libvirts.CreatePool(libvirts.ToDomainPool(name), path)
	if err != nil {
		return nil, err
	}
	vol, err := libvirts.CreateBackingVolume(pol.Name, disk, backingFle, backingFmt)
	if err != nil {
		return nil, err
	}
	return vol.GetXMLObj()
}

// name: Domain name.
func RemovePool(name string) error {
	pol := &libvirts.Pool{Name: name}
	if err := pol.Remove(); err != nil {
		return err
	}
	return nil
}

func CleanPool(name string) error {
	pol := &libvirts.Pool{Name: name}
	if err := pol.Clean(); err != nil {
		return err
	}
	return nil
}

var pubDir = ""

func SetStatic(dir string) {
	pubDir = dir
}

func GetStatic() string {
	return pubDir
}

func GetFile(name string) string {
	return fmt.Sprintf("%s/%s", pubDir, name)
}

func ParseFiles(w http.ResponseWriter, name string, data interface{}) error {
	file := path.Base(name)
	tmpl, err := template.New(file).Funcs(template.FuncMap{
		"prettyBytes":  libstar.PrettyBytes,
		"prettyKBytes": libstar.PrettyKBytes,
		"prettySecs":   libstar.PrettySecs,
		"prettyPCI":    libstar.PrettyPCI,
		"prettyDrive":  libstar.PrettyDrive,
	}).ParseFiles(name)
	if err != nil {
		_, _ = fmt.Fprintf(w, "template.ParseFiles %s", err)
		return err
	}
	if err := tmpl.Execute(w, data); err != nil {
		_, _ = fmt.Fprintf(w, "template.ParseFiles %s", err)
		return err
	}
	return nil
}

func GetAuth(req *http.Request) (name, pass string, ok bool) {
	if t, err := req.Cookie("session-id"); err == nil {
		if sess := service.SERVICE.Session.Get(t.Value); sess != nil {
			now := time.Now()
			expired, _ := libstar.GetLocalTime(time.RFC3339, sess.Expires)
			if now.Before(expired) {
				name, pass, ok = ParseBasicAuth(sess.Value)
			}
		}
	} else {
		name, pass, ok = req.BasicAuth()
	}
	return name, pass, ok
}

func GetUser(req *http.Request) (schema.User, bool) {
	name, _, _ := GetAuth(req)
	libstar.Debug("GetUser %s", name)
	return service.SERVICE.Users.Get(name)
}

func ParseBasicAuth(auth string) (username, password string, ok bool) {
	c, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func UpdateCookie(w http.ResponseWriter, req *http.Request, value string) {
	uuid := ""
	if obj, err := req.Cookie("session-id"); err == nil {
		uuid = obj.Value
	} else {
		uuid = libstar.GenToken(32)
	}
	expired := time.Now().Add(time.Minute * 15)
	sess := &schema.Session{
		Uuid:    uuid,
		Client:  req.RemoteAddr,
		Value:   base64.StdEncoding.EncodeToString([]byte(value)),
		Expires: expired.Format(time.RFC3339),
	}
	service.SERVICE.Session.Add(sess)
	http.SetCookie(w, &http.Cookie{
		Name:    "session-id",
		Value:   uuid,
		Path:    "/",
		Expires: expired,
	})
}
