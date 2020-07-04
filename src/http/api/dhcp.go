package api

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/danieldin95/lightstar/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type DHCPLease struct {
}

func (l DHCPLease) Router(router *mux.Router) {
	router.HandleFunc("/api/network/all/lease", l.GET).Methods("GET")
	router.HandleFunc("/api/network/{id}/lease", l.GET).Methods("GET")
}

func (l DHCPLease) Get(data schema.DHCPLeases) error {
	leases, err := libvirtn.ListLeases()
	if err != nil {
		return err
	}
	for addr, l := range leases {
		data[addr] = schema.DHCPLease{
			Mac:      l.Mac,
			IPAddr:   l.IPAddr,
			Prefix:   l.Prefix,
			Hostname: l.Hostname,
			Type:     l.Type,
		}
	}
	return nil
}

func (l DHCPLease) GET(w http.ResponseWriter, r *http.Request) {
	uuid, ok := GetArg(r, "id")
	if !ok || uuid == "all" {
		data := make(schema.DHCPLeases, 128)
		if err := l.Get(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ResponseJson(w, data)
	} else {
		leases, err := libvirtn.LookupLeases(uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		list := schema.ListDHCPLease{
			Items: make([]schema.DHCPLease, 0, 32),
		}
		for _, l := range leases {
			list.Items = append(list.Items, schema.DHCPLease{
				Mac:      l.Mac,
				IPAddr:   l.IPAddr,
				Prefix:   l.Prefix,
				Hostname: l.Hostname,
				Type:     l.Type,
			})
		}
		ResponseJson(w, list)
	}
}

func (l DHCPLease) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (l DHCPLease) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (l DHCPLease) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}
