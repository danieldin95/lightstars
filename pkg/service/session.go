package service

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
	"sync"
	"time"
)

type Session struct {
	Lock    sync.RWMutex
	File    string
	Session map[string]*schema.Session `json:"session"`
}

func (s *Session) Save() error {
	s.Lock.RLock()
	sess := make(map[string]*schema.Session, 32)
	for key, obj := range s.Session {
		now := time.Now()
		expired, _ := libstar.GetLocalTime(time.RFC3339, obj.Expires)
		if now.After(expired) {
			continue
		}
		sess[key] = obj
	}
	s.Lock.RUnlock()
	if err := libstar.JSON.MarshalSave(&sess, s.File, false); err != nil {
		return err
	}
	return nil
}

func (s *Session) Load(file string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.File = file
	if err := libstar.JSON.UnmarshalLoad(&s.Session, file); err != nil {
		return err
	}
	return nil
}

func (s *Session) Add(obj *schema.Session) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Session[obj.Uuid] = obj
}

func (s *Session) Mod(obj *schema.Session) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if _, ok := s.Session[obj.Uuid]; ok {
		s.Session[obj.Uuid] = obj
	}
}

func (s *Session) Del(key string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if _, ok := s.Session[key]; ok {
		delete(s.Session, key)
	}
}

func (s *Session) Get(key string) *schema.Session {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.Session[key]
}

func (s *Session) AddAndSave(obj *schema.Session) {
	s.Add(obj)
	if err := s.Save(); err != nil {
		libstar.Warn("Session.AddAndSave %s", err)
	}
}

func (s *Session) List(user string) <-chan *schema.Session {
	c := make(chan *schema.Session, 128)
	go func() {
		s.Lock.RLock()
		defer s.Lock.RUnlock()

		for _, h := range s.Session {
			c <- h
		}
		c <- nil //Finish channel by nil.
	}()
	return c
}
