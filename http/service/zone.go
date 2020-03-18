package service

import (
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"sync"
)

type Zone struct {
	Lock sync.RWMutex
	Name string
	Host map[string]*schema.Host
}

func (l *Zone) Load(file string) error {
	if err := libstar.JSON.UnmarshalLoad(&l.Host, file); err != nil {
		return err
	}
	libstar.Debug("Zone.Load %v", l.Host)
	for name, host := range l.Host {
		if host == nil {
			continue
		}
		host.Name = name
		libstar.Debug("Zone.Load %s", host)
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
