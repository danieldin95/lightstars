package client

import (
	"fmt"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"testing"
)

func TestDHCPLease_Get(t *testing.T) {
	api := DHCPLease{
		Client: Client{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: "admin:123",
			},
		},
		Host: "https://localhost:10080",
	}
	les := map[string]schema.DHCPLease{}
	fmt.Println(api.Get(&les), les)
	les = map[string]schema.DHCPLease{}
	api.Client.Auth.Username = "123"
	fmt.Println(api.Get(&les), les)
}

func TestProxyTcp_Get(t *testing.T) {
	api := ProxyTcp{
		Client: Client{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: "admin:123",
			},
		},
		Host: "https://localhost:10080",
	}
	var ps []schema.Target
	fmt.Println(api.Get(&ps), ps)
}

func TestInstance_Get(t *testing.T) {
	api := Instance{
		Client: Client{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: "admin:123",
			},
		},
		Host: "https://localhost:10080",
	}
	var data schema.List
	fmt.Println(api.Get(&data), data)
}
