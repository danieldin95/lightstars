package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/http/service"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"text/template"
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
		w.Write(str)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResponseXML(w http.ResponseWriter, v string) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(v))
}

func ResponseMsg(w http.ResponseWriter, code int, message string) {
	ret := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	}
	ResponseJson(w, ret)
}

func Slot2Disk(slot uint8) string {
	return libvirtc.DISK.Slot2DiskName(slot)
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
func DelVolumeAndPool(name string) error {
	err := libvirts.RemovePool(libvirts.ToDomainPool(name))
	if err != nil {
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
		fmt.Fprintf(w, "template.ParseFiles %s", err)
		return err
	}
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Fprintf(w, "template.ParseFiles %s", err)
		return err
	}
	return nil
}

func GetAuth(req *http.Request) (name, pass string, ok bool) {
	if t, err := req.Cookie("token"); err == nil {
		name, pass, ok = ParseBasicAuth(t.Value)
	} else {
		name, pass, ok = req.BasicAuth()
	}
	return name, pass, ok
}

func GetUser(req *http.Request) (schema.User, bool) {
	name, _, _ := GetAuth(req)
	libstar.Debug("GetUser %s", name)
	return service.USERS.Get(name)
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
