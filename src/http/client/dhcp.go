package client

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
)

type DHCPLease struct {
	Client
}

func (api DHCPLease) Url() string {
	return api.Host + "/api/network/all/lease"
}

func (api DHCPLease) Get(data *schema.DHCPLeases) error {
	client := api.NewRequest(api.Url())
	if err := api.GetJSON(client, data); err != nil {
		libstar.Error("DHCPLease.Get %s", err)
		return err
	}
	return nil
}
