package main

import (
	"fmt"
	"github.com/danieldin95/lightstar/http"
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
	http  := http.NewServer("0.0.0.0:10080")
	go http.Start()

	Wait()
}