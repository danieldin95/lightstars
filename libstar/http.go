package libstar

import (
	"crypto/tls"
	"io"
	"net/http"
)

type HttpClient struct {
	Method    string
	Url       string
	Payload   io.Reader
	Auth      Auth
	TlsConfig *tls.Config
}

func (cl *HttpClient) Do() (*http.Response, error) {
	if cl.Method == "" {
		cl.Method = "GET"
	}
	if cl.TlsConfig == nil {
		cl.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	}
	req, err := http.NewRequest(cl.Method, cl.Url, cl.Payload)
	if err != nil {
		return nil, err
	}
	if cl.Auth.Type == "basic" {
		req.SetBasicAuth(cl.Auth.Username, cl.Auth.Password)
	}
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cl.TlsConfig,
		},
	}
	return client.Do(req)
}
