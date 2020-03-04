package main

import (
	"flag"
	"fmt"
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"os"
	"os/signal"
	"syscall"
)

func Wait() {
	x := make(chan os.Signal)
	signal.Notify(x, os.Interrupt, syscall.SIGTERM)
	signal.Notify(x, os.Interrupt, syscall.SIGKILL)
	signal.Notify(x, os.Interrupt, syscall.SIGQUIT) //CTL+/
	signal.Notify(x, os.Interrupt, syscall.SIGINT)  //CTL+C

	<-x
	fmt.Println("Done")
}

func main() {
	staticDir := "static"
	crtDir := "ca"
	authFile := ".auth"
	listen := "0.0.0.0:10080"
	hyper := "qemu:///system"
	verbose := 2
	logFile := "/var/log/lightstar.log"

	flag.StringVar(&listen, "listen", listen, "the address http listen.")
	flag.IntVar(&verbose, "log:level", verbose, "logger level")
	flag.StringVar(&hyper, "hyper", hyper, "hypervisor connecting to.")
	flag.StringVar(&crtDir, "crt:dir", crtDir, "he directory X509 certificate file on.")
	flag.StringVar(&staticDir, "static:dir", staticDir, "the directory to serve files from.")
	flag.StringVar(&authFile, "auth:file", authFile, "the file saved administrator auth")
	flag.Parse()

	libstar.Init(logFile, verbose)
	// initialize storage
	libvirts.DATASTOR.Init()

	libvirtc.SetHyper(hyper)
	libvirts.SetHyper(hyper)
	libvirtn.SetHyper(hyper)

	h := http.NewServer(listen, staticDir, authFile)
	if _, err := os.Stat(crtDir); !os.IsNotExist(err) {
		h.SetCert(crtDir+"/private.key", crtDir+"/crt.pem")
	}

	go h.Start()

	Wait()
}
