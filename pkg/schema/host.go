package schema

import "net/url"

type Host struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Schema   string `json:"schema"`
	Port     string `json:"port"`
}

func (h *Host) Initialize() {
	if url, err := url.Parse(h.Url); err == nil {
		h.Port = url.Port()
		h.Schema = url.Scheme
		h.Hostname = url.Hostname()
	}
}
