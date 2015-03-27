package teamgress

import (
	"fmt"
)

type User struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (u *User) String() string {
	return fmt.Sprintf("\"%s\" (%s)", u.Name, u.Key)
}

func (c *Conf) _indexByUsername() {
	for _, u := range c.Users {
		c.indexUsername[u.Key] = u
	}
}

func (c *Conf) GetByUsername(username string) User {
	if u, found := c.indexUsername[username]; found {
		return *u
	}

	return User{}
}
