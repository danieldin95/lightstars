package main

import (
	"flag"
	"github.com/danieldin95/lightstar/src/compute/libvirtc"
	"github.com/danieldin95/lightstar/src/http"
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/network/libvirtn"
	"github.com/danieldin95/lightstar/src/service"
	"github.com/danieldin95/lightstar/src/storage"
	"github.com/danieldin95/lightstar/src/storage/libvirts"
	"os"
)

type StarConfig struct {
	StaticDir string `json:"dir.static"`
	CrtDir    string `json:"dir.crt"`
	ConfDir   string `json:"-"`
	Hyper     string `json:"hyper"`
	Verbose   int    `json:"log.level"`
	LogFile   string `json:"log.file"`
	Listen    string `json:"listen"`
}

var cfg = StarConfig{
	StaticDir: "static",
	CrtDir:    "ca",
	ConfDir:   "/etc/lightstar",
	Listen:    "0.0.0.0:10080",
	Hyper:     "qemu:///system",
	LogFile:   "/var/log/lightstar.log",
	Verbose:   2,
}

func pprof(file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return
	}
	addr := ""
	if err := libstar.JSON.UnmarshalLoad(&addr, file); err != nil {
		libstar.Warn("pprof.JSON.UnmarshalLoad %v", err)
	}
	p := &libstar.PProf{
		Listen: addr,
	}
	p.Start()
}

func main() {
	flag.StringVar(&cfg.Listen, "listen", cfg.Listen, "the address http listen.")
	flag.IntVar(&cfg.Verbose, "log:level", cfg.Verbose, "logger level")
	flag.StringVar(&cfg.Hyper, "hyper", cfg.Hyper, "hypervisor connecting to.")
	flag.StringVar(&cfg.CrtDir, "crt:dir", cfg.CrtDir, "the directory X509 certificate file on.")
	flag.StringVar(&cfg.StaticDir, "static:dir", cfg.StaticDir, "the directory to serve files from.")
	flag.StringVar(&cfg.ConfDir, "conf", cfg.ConfDir, "the directory configuration on")
	flag.Parse()

	libstar.PreNotify()
	// Check and Start pprof.
	pprof(cfg.ConfDir + "/pprof.json")
	libstar.Init(cfg.LogFile, cfg.Verbose)
	// Initialize storage
	storage.DATASTOR.Init()
	service.SERVICE.Load(cfg.ConfDir)
	// Initialize hyper
	_, _ = libvirtc.SetHyper(cfg.Hyper)
	_, _ = libvirts.SetHyper(cfg.Hyper)
	_, _ = libvirtn.SetHyper(cfg.Hyper)
	// Configure cert and auth.
	authFile := cfg.ConfDir + "/auth.json"
	h := http.NewServer(cfg.Listen, cfg.StaticDir, authFile)
	if _, err := os.Stat(cfg.CrtDir); !os.IsNotExist(err) {
		h.SetCert(cfg.CrtDir+"/private.key", cfg.CrtDir+"/crt.pem")
	}
	// Start
	go h.Start()
	libstar.SdNotify()

	libstar.Wait()
}
