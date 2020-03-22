package client

import "github.com/danieldin95/lightstar/libstar"

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
