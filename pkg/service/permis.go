package service

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"net/http"
	"strings"
)

type RuleValue struct {
	Value string
	Type  string // argument or path
}

type Rule struct {
	Path   string      `json:"path"`
	Values []RuleValue `json:"-"`
	Type   string      `json:"type"`   // prefix, suffix or path
	Method string      `json:"method"` // get, post, put, delete
	Action string      `json:"action"` // permit or deny
}

type RouteMatcher struct {
	Match []Rule
}

func (r *RouteMatcher) Add(m Rule) {
	if m.Values == nil {
		m.Values = make([]RuleValue, 0, 32)
	}
	for _, v := range strings.Split(m.Path, "/") {
		ri := RuleValue{
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

type Permission struct {
	Guest RouteMatcher `json:"guest"`
}

func (p *Permission) Load(file string) error {
	rules := struct {
		Guest []Rule `json:"guest"`
	}{}
	if err := libstar.JSON.UnmarshalLoad(&rules, file); err != nil {
		return err
	}
	// guest permission.
	for _, v := range rules.Guest {
		p.Guest.Add(v)
	}
	return nil
}

func (p *Permission) Has(req *http.Request) bool {
	path := req.URL.Path
	values := strings.Split(path, "/")
	for _, k := range p.Guest.Match {
		if req.Method != strings.ToUpper(k.Method) {
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
