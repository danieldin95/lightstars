package main

import (
	"flag"
	"fmt"
	"github.com/danieldin95/lightstar/http"
	"github.com/danieldin95/lightstar/libstar"
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

var (
	Version string
	Date    string
	Commit  string
)

func init() {
	libstar.Info("version is %s", Version)
	libstar.Info("built on %s", Date)
	libstar.Info("commit at %s", Commit)
}

func main() {
	var dir string
	var listen string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from.")
	flag.StringVar(&listen, "listen", "0.0.0.0:10080", "the address http listen.")
	flag.Parse()

	http  := http.NewServer(listen, dir)
	go http.Start()

	Wait()
}