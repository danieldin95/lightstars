package main

import (
	"flag"
	"fmt"
	"github.com/danieldin95/lightstar/http/client"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/proxy"
	"github.com/danieldin95/lightstar/schema"
	"strings"
)

func GetPorts(host, auth string) []schema.Target {
	api := client.ProxyTcp{
		Client: client.Client{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: auth,
			},
		},
		Host: host,
	}
	return api.Get()
}

type PixConfig struct {
	Url     string   `json:"url"`
	Auth    string   `json:"Auth"`
	Target  []string `json:"target"`
	Verbose int      `json:"log.verbose"`
	LogFile string   `json:"log.file"`

	Targets []schema.Target
}

func (cfg *PixConfig) Parse() *PixConfig {
	tgt := ""
	file := "lightpix.json"
	cfg.Url = "https://localhost:10080"
	cfg.Auth = "admin:123456"
	cfg.Verbose = 2
	cfg.LogFile = "lightpix.log"

	flag.StringVar(&file, "conf", file, "the configuration file")
	flag.StringVar(&cfg.Url, "url", cfg.Url, "the url path.")
	flag.StringVar(&tgt, "tgt", tgt, "target list by comma.")
	flag.StringVar(&cfg.Auth, "auth", cfg.Auth, "the auth login to.")
	flag.IntVar(&cfg.Verbose, "log:level", cfg.Verbose, "logger level")
	flag.Parse()

	cfg.Target = strings.Split(tgt, ",")
	if err := libstar.JSON.UnmarshalLoad(&cfg, file); err != nil {
		libstar.Warn("main %s", err)
	}
	if cfg.Targets == nil {
		cfg.Targets = make([]schema.Target, 0, 32)
	}
	for _, t := range cfg.Target {
		cfg.Targets = append(cfg.Targets, schema.Target{
			Name:   "custom",
			Target: t,
		})
	}
	return cfg
}

func main() {
	cfg := &PixConfig{}
	cfg.Parse()
	libstar.Init(cfg.LogFile, cfg.Verbose)

	ports := GetPorts(cfg.Url, cfg.Auth)
	for _, port := range ports {
		cfg.Targets = append(cfg.Targets, port)
	}
	pri := proxy.Proxy{
		Target: cfg.Targets,
		Client: &proxy.WsClient{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: cfg.Auth,
			},
			Url: cfg.Url + "/ext/tcpsocket?target=",
		},
	}
	pri.Initialize()
	pri.Start()
	go func() {
		for {
			input := ""
			fmt.Scanln(&input)
			for _, tgt := range pri.Target {
				if l, ok := pri.Listen[tgt.Target]; ok {
					libstar.Info("main %-15s %-20s on %-15s", l.Tgt.Name, l.Tgt.Target, l.Listen)
				}
			}
		}
	}()
	defer pri.Stop()
	libstar.Wait()
}
