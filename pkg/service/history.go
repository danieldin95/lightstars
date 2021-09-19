package service

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
	"sync"
)

type History struct {
	Lock    sync.RWMutex
	File    string
	History []*schema.History `json:"history"`
}

func (h *History) Save() error {
	h.Lock.RLock()
	defer h.Lock.RUnlock()
	if err := libstar.JSON.MarshalSave(&h.History, h.File, false); err != nil {
		return err
	}
	return nil
}

func (h *History) Load(file string) error {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.File = file
	if err := libstar.JSON.UnmarshalLoad(&h.History, file); err != nil {
		return err
	}
	return nil
}

func (h *History) Add(obj *schema.History) {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	h.History = append(h.History, obj)
}

func (h *History) List(user string) <-chan *schema.History {
	c := make(chan *schema.History, 128)
	go func() {
		h.Lock.RLock()
		defer h.Lock.RUnlock()

		for _, h := range h.History {
			if user == "" || user == h.User {
				c <- h
			}
		}
		c <- nil //Finish channel by nil.
	}()
	return c
}
