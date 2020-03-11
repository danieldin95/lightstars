package api

import (
	"github.com/danieldin95/lightstar/libstar"
	"net/http"
	"strings"
)

type RouteMatch struct {
	Path   string
	Values []string
	Type   string // prefix or path
	Method string
	Action string // permit or deny
}

type RouteMatcher struct {
	Match []RouteMatch
}

var ROUTER = &RouteMatcher{
	Match: make([]RouteMatch, 0, 32),
}

func (r *RouteMatcher) Add(m RouteMatch) {
	libstar.Info("RouteMatcher.Add %v", m)
	m.Values = strings.Split(m.Path, "/")
	r.Match = append(r.Match, m)
}

func HasPermission(req *http.Request) bool {
	user, _ := GetUser(req)
	if user.Type == "admin" {
		return true
	}

	path := req.URL.Path
	values := strings.Split(path, "/")
	for _, k := range ROUTER.Match {
		if req.Method != k.Method {
			continue
		}
		switch k.Type {
		case "prefix":
			if strings.HasPrefix(path, k.Path) {
				return true
			}
		case "suffix":
			if strings.HasSuffix(path, k.Path) {
				return true
			}
		default: // strict path matching.
			if len(k.Values) != len(values) {
				continue
			}
			matched := true
			for i, v := range k.Values {
				if values[i] == v {
					continue
				}
				if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
					continue
				}
				matched = false
			}
			if !matched {
				continue
			}
			if k.Action == "permit" {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func init() {
	// guest permission.
	ROUTER.Add(RouteMatch{
		Path:   "/",
		Type:   "prefix",
		Method: "GET",
		Action: "permit",
	})
	ROUTER.Add(RouteMatch{
		Path:   "/api/instance/{id}",
		Method: "PUT",
		Action: "permit",
	})
}
