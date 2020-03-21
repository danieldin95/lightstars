package main

import (
	"flag"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/proxy"
	"strings"
)

func main() {
	url := "wss://localhost:10080"
	auth := "admin:123456"
	verbose := 2
	tgt := "localhost:22"

	flag.StringVar(&url, "url", url, "the url path.")
	flag.StringVar(&tgt, "tgt", tgt, "the target proxied to.")
	flag.StringVar(&auth, "auth", auth, "the auth login to.")
	flag.IntVar(&verbose, "log:level", verbose, "logger level")
	flag.Parse()

	name := auth
	password := ""
	if strings.Contains(auth, ":") {
		name = strings.SplitN(auth, ":", 2)[0]
		password = strings.SplitN(auth, ":", 2)[1]
	}
	pri := proxy.Proxy{
		Target: strings.Split(tgt, ","),
		Client: &proxy.WsClient{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: name,
				Password: password,
			},
			Url: url + "/ext/tcpsocket?target=",
		},
	}
	pri.Initialize().Start()
	defer pri.Stop()
	libstar.Wait()
}
