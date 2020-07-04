package http

import (
	"sync"
	"time"
)

type session struct {
	lock    sync.RWMutex
	session map[string]struct {
		UUID   string
		Expire time.Time
		Client string
	}
	done   chan bool
	ticker *time.Ticker
}

func (s *session) Update(uuid string, expire time.Time) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if item, ok := s.session[uuid]; ok {
		item.Expire = expire
	}
}

func (s *session) Add(client, uuid string, expire time.Time) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.session[uuid] = struct {
		UUID   string
		Expire time.Time
		Client string
	}{UUID: uuid, Expire: expire, Client: client}
}

func (s *session) Del(uuid string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.session[uuid]; ok {
		delete(s.session, uuid)
	}
}

func (s *session) Expired(uuid string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if item, ok := s.session[uuid]; ok {
		return time.Now().After(item.Expire)
	}
	return true
}

func (s *session) loop() {
	expired := make([]string, 32)
	for uuid, item := range s.session {
		if time.Now().After(item.Expire) {
			expired = append(expired, uuid)
		}
	}
	for _, uuid := range expired {
		delete(s.session, uuid)
	}
}

func (s *session) Start() {
	for {
		select {
		case <-s.done:
			return
		case <-s.ticker.C:
			s.lock.Lock()
			s.loop()
			s.lock.Unlock()
		}
	}
}

func (s *session) Stop() {
	s.done <- true
}

var Session = session{
	session: make(map[string]struct {
		UUID   string
		Expire time.Time
		Client string
	}, 1024),
	done:   make(chan bool, 2),
	ticker: time.NewTicker(5 * time.Second),
}
