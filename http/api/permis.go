package api

import (
	"github.com/danieldin95/lightstar/libstar"
	"net/http"
	"strings"
)

type RouteItem struct {
	Value string
	Type  string // argument or path
}

type RouteMatch struct {
	Path   string
	Values []RouteItem
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
	if m.Values == nil {
		m.Values = make([]RouteItem, 0, 32)
	}
	for _, v := range strings.Split(m.Path, "/") {
		ri := RouteItem{
			Value: v,
			Type:  "path",
		}
		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			ri.Type = "argument"
		}
		m.Values = append(m.Values, ri)
	}
	libstar.Debug("RouteMatcher.Add %v", m)
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
				if values[i] == v.Value || v.Type == "argument" {
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
		Path:   "/login",
		Type:   "prefix",
		Method: "POST",
		Action: "permit",
	})
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
