package client

import (
	"encoding/json"
	"github.com/danieldin95/lightstar/libstar"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Auth libstar.Auth
}

func (cl Client) NewRequest(url string) *libstar.HttpClient {
	client := &libstar.HttpClient{
		Auth: libstar.Auth{
			Type:     "basic",
			Username: cl.Auth.Username,
			Password: cl.Auth.Password,
		},
		Url: url,
	}
	return client
}

func (cl Client) GetJSON(client *libstar.HttpClient, v interface{}) error {
	r, err := client.Do()
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return libstar.NewErr(r.Status)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	libstar.Debug("Client.GetJSON %s", body)
	if err := json.Unmarshal([]byte(body), v); err != nil {
		return err
	}
	return nil
}
