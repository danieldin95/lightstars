package client

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
)

type Instance struct {
	Client
	Host string
	Name string
}

func (api Instance) Url() string {
	if api.Name == "" {
		return api.Host + "/api/instance"
	}
	return api.Host + "/api/instance/" + api.Name
}

func (api Instance) Get(data *schema.List) error {
	client := api.NewRequest(api.Url())
	if err := api.GetJSON(client, data); err != nil {
		libstar.Error("Instance.Get %s", err)
		return err
	}
	return nil
}
