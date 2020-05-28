package libstar

import (
	"github.com/coreos/go-systemd/v22/daemon"
	"os"
	"strconv"
)

func PreNotify() {
}

func SdNotify() {
	go daemon.SdNotify(false, daemon.SdNotifyReady)
}

func SavePID(file string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	pid := os.Getpid()
	if _, err := f.Write([]byte(strconv.Itoa(pid))); err != nil {
		return err
	}
	return nil
}
