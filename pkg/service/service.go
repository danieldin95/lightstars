package service

import (
	"github.com/danieldin95/lightstar/pkg/libstar"
	"github.com/danieldin95/lightstar/pkg/schema"
	"time"
)

type Service struct {
	Zone       *Zone
	Users      *Users
	Permission *Permission
	History    *History
	Session    *Session
	Done       chan bool
	Ticker     *time.Ticker
}

var SERVICE = Service{
	Zone: &Zone{
		Host: make(map[string]*schema.Host, 32),
	},
	Users: &Users{
		Users: make(map[string]*schema.User, 32),
	},
	Permission: &Permission{
		Guest: RouteMatcher{
			Match: make([]Rule, 0, 32),
		},
	},
	History: &History{
		History: make([]*schema.History, 0, 128),
	},
	Session: &Session{
		Session: make(map[string]*schema.Session, 12),
	},
	Done:   make(chan bool),
	Ticker: time.NewTicker(2 * time.Second),
}

func (s *Service) Load(path string) {
	if err := s.Zone.Load(path + "/zone.json"); err != nil {
		libstar.Error("Service.Load.Zone %s", err)
	}
	if err := s.Users.Load(path + "/auth.json"); err != nil {
		libstar.Error("Service.Load.Users %s", err)
	}
	if err := s.Permission.Load(path + "/permission.json"); err != nil {
		libstar.Error("Service.Load.Permission %s", err)
	}
	if err := s.History.Load(path + "/history.json"); err != nil {
		libstar.Error("Service.Load.History %s", err)
	}
	if err := s.Session.Load(path + "/session.json"); err != nil {
		libstar.Error("Service.Load.Session %s", err)
	}
	libstar.Debug("Service.Load %v", s)
}

func (s *Service) Flush() {
	_ = s.Session.Save()
	_ = s.History.Save()
}

func (s *Service) Loop() {
	for {
		select {
		case <-s.Done:
			return
		case <-s.Ticker.C:
			s.Flush()
		}
	}
}
