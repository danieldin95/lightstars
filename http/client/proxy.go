package client

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
)

type ProxyTcp struct {
	Client
	Host string
}

func (pro ProxyTcp) Url() string {
	return pro.Host + "/api/proxy/tcp"
}

func (pro ProxyTcp) Get() []schema.Target {
	ports := make([]schema.Target, 0, 32)
	client := pro.NewRequest(pro.Url())
	r, err := client.Do()
	if err == nil {
		libstar.GetJSON(r.Body, &ports)
	} else {
		libstar.Error("ProxyTcp.Get %s", err)
	}
	libstar.Debug("ProxyTcp.Get %s", ports)
	return ports
}
