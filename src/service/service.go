package service

import (
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
)

type Service struct {
	Zone       Zone
	Users      Users
	Permission Permission
}

var SERVICE = Service{
	Zone: Zone{
		Host: make(map[string]*schema.Host, 32),
	},
	Users: Users{
		Users: make(map[string]*schema.User, 32),
	},
	Permission: Permission{
		Guest: RouteMatcher{
			Match: make([]Rule, 0, 32),
		},
	},
}

func (s *Service) Load(path string) {
	if err := s.Zone.Load(path + "/zone.json"); err != nil {
		libstar.Error("Service.Load.Zone %s", err)
	}
	libstar.Debug("Service.Load %v", s.Zone)
	if err := s.Users.Load(path + "/auth.json"); err != nil {
		libstar.Error("Service.Load.Users %s", err)
	}
	libstar.Debug("Service.Load %v", s.Users)
	if err := s.Permission.Load(path + "/permission.json"); err != nil {
		libstar.Error("Service.Load.Permission %s", err)
	}
	libstar.Debug("Service.Load %v", s.Permission)
}
