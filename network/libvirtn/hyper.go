package libvirtn

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/libvirt/libvirt-go"
)

type HyperVisor struct {
	Name string
	Conn *libvirt.Connect
}

func (h *HyperVisor) Open() error {
	if hyper.Conn != nil {
		if _, err := hyper.Conn.GetVersion(); err != nil {
			libstar.Error("HyperVisor.Open %s", err)
			hyper.Conn.Close()
			hyper.Conn = nil
		}
	}
	if hyper.Conn == nil {
		conn, err := libvirt.NewConnect(hyper.Name)
		if err != nil {
			return err
		}
		hyper.Conn = conn
	}
	if hyper.Conn == nil {
		return libstar.NewErr("Not connect.")
	}
	return nil
}

var hyper = HyperVisor{
	Name: "qemu:///system",
}

func GetHyper() (*HyperVisor, error) {
	return &hyper, hyper.Open()
}

func SetHyper(name string) (*HyperVisor, error) {
	if name == hyper.Name {
		return &hyper, nil
	}
	hyper.Name = name

	conn, err := libvirt.NewConnect(hyper.Name)
	if err != nil {
		return &hyper, err
	}
	hyper.Conn = conn
	return &hyper, nil
}

func CloseHyper(name string) {
	hyper.Conn.Close()
}
