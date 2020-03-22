package client

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
)

type DHCPLease struct {
	Client
}

func (api DHCPLease) Url() string {
	return api.Host + "/api/dhcp/lease"
}

func (api DHCPLease) Get(data *schema.ListLeases) error {
	client := api.NewRequest(api.Url())
	if err := api.GetJSON(client, data); err != nil {
		libstar.Error("DHCPLease.Get %s", err)
		return err
	}
	return nil
}
