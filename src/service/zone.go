package service

import (
	"github.com/danieldin95/lightstar/src/libstar"
	"github.com/danieldin95/lightstar/src/schema"
	"sync"
)

type Zone struct {
	Lock sync.RWMutex
	Name string
	Host map[string]*schema.Host
}

func (l *Zone) Load(file string) error {
	var hosts map[string]*schema.Host

	if l.Host == nil {
		l.Host = make(map[string]*schema.Host, 32)
	}
	l.Add(&schema.Host{Name: "default", Url: ""})
	if err := libstar.JSON.UnmarshalLoad(&hosts, file); err != nil {
		return err
	}
	libstar.Debug("Zone.Load %v", hosts)
	for name, host := range hosts {
		if host == nil {
			continue
		}
		host.Name = name
		host.Initialize()
		libstar.Debug("Zone.Load %v", host)
		l.Add(host)
	}
	return nil
}

func (l *Zone) Get(name string) *schema.Host {
	l.Lock.RLock()
	defer l.Lock.RUnlock()
	return l.Host[name]
}

func (l *Zone) Add(h *schema.Host) {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	l.Host[h.Name] = h
}

func (l *Zone) List() <-chan *schema.Host {
	c := make(chan *schema.Host, 128)
	go func() {
		l.Lock.RLock()
		defer l.Lock.RUnlock()

		for _, h := range l.Host {
			c <- h
		}
		c <- nil //Finish channel by nil.
	}()
	return c
}
