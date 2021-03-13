package api

import (
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/schema"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"time"
)

type Snapshot struct {
}

func (in Snapshot) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/snapshot", in.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/snapshot", in.Post).Methods("POST")
	router.HandleFunc("/api/instance/{id}/snapshot/{name}", in.Get).Methods("GET")
	router.HandleFunc("/api/instance/{id}/snapshot/{name}/revert", in.Revert).Methods("PUT")
	router.HandleFunc("/api/instance/{id}/snapshot/{name}", in.Delete).Methods("DELETE")
}

func (in Snapshot) Get(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer dom.Free()

	list := schema.ListSnapshot{
		Items: make([]schema.Snapshot, 0, 32),
	}
	name, ok := GetArg(r, "name")
	if !ok {
		sns, err := dom.ListAllSnapshots(0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, obj := range sns {
			if snx := libvirtc.NewSnapshotXMLFromDom(&obj); snx != nil {
				list.Items = append(list.Items, schema.Snapshot{
					Name:   snx.Name,
					State:  snx.State,
					Uptime: time.Now().Unix() - snx.CreateAt,
				})
			}
			_ = obj.Free()
		}
	} else {
		if obj, err := dom.SnapshotLookupByName(name, 0); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if snx := libvirtc.NewSnapshotXMLFromDom(obj); snx != nil {
			list.Items = append(list.Items, schema.Snapshot{
				Name:   snx.Name,
				State:  snx.State,
				Uptime: time.Now().Unix() - snx.CreateAt,
			})
			_ = obj.Free()
		}
	}
	sort.SliceStable(list.Items, func(i, j int) bool {
		return list.Items[i].Uptime < list.Items[j].Uptime
	})
	list.Metadata.Size = len(list.Items)
	list.Metadata.Total = len(list.Items)
	ResponseJson(w, list)
}

func (in Snapshot) Post(w http.ResponseWriter, r *http.Request) {
	conf := &schema.Snapshot{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	xml := libvirtc.SnapshotXML{
		Name: conf.Name,
	}
	if snx, err := dom.CreateSnapshotXML(xml.Encode(), 0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		snx.Free()
	}
	ResponseMsg(w, 0, "success")
}

func (in Snapshot) Revert(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	name, ok := GetArg(r, "name")
	if ok {
		if obj, err := dom.SnapshotLookupByName(name, 0); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if err := obj.RevertToSnapshot(0); err != nil {
				_ = obj.Free()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = obj.Free()
		}
	}
	ResponseMsg(w, 0, "success")
}

func (in Snapshot) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	name, ok := GetArg(r, "name")
	if ok {
		if obj, err := dom.SnapshotLookupByName(name, 0); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if err := obj.Delete(0); err != nil {
				_ = obj.Free()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = obj.Free()
		}
	}
	ResponseMsg(w, 0, "success")
}
