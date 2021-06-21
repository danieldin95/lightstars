package client

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
)

type ProxyTcp struct {
	Client
}

func (api ProxyTcp) Url() string {
	return api.Host + "/api/proxy/tcp"
}

func (api ProxyTcp) Get(data *[]schema.Target) error {
	client := api.NewRequest(api.Url())
	if err := api.GetJSON(client, data); err != nil {
		libstar.Error("DHCPLease.Get %s", err)
		return err
	}
	return nil
}
