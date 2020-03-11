package service

import (
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"sync"
)

type Users struct {
	lock  sync.RWMutex
	file  string
	users map[string]schema.User `json:"user"`
}

func (u Users) Save() error {
	u.lock.RLock()
	defer u.lock.RUnlock()

	if err := libstar.JSON.MarshalSave(&u.users, u.file, true); err != nil {
		libstar.Error("Server.LoadToken: %s", err)
		return err
	}
	return nil
}

func (u Users) Load(file string) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.file = file
	if err := libstar.JSON.UnmarshalLoad(&u.users, file); err != nil {
		libstar.Error("Users.Load: %s", err)
	}
	token := ""
	for k, v := range u.users {
		if k == "admin" {
			token = v.Password
		}
	}
	if token == "" {
		token = libstar.GenToken(32)
		u.users["admin"] = schema.User{
			Type:     "admin",
			Name:     "admin",
			Password: token,
		}
		libstar.JSON.MarshalSave(&u.users, u.file, true)
	}
}

func (u Users) Get(name string) (schema.User, bool) {
	u.lock.RLock()
	defer u.lock.RUnlock()

	user, ok := u.users[name]
	return user, ok
}

var USERS = Users{
	users: make(map[string]schema.User, 32),
}
