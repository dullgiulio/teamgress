package main

type User struct {
	Name      string `json:"name"`
	UnixName  string `json:"unix"`
	OtherName string `json:"other"`
	Avatar    string `json:"avatar"`
}

func (u *User) String() string {
	return u.Name
}

func (c *conf) _indexByUsername() {
	for _, u := range c.Users {
		c.indexUsername[u.UnixName] = u
	}
}

func (c *conf) getByUsername(username string) User {
	c.mux.Lock()
	defer c.mux.Unlock()

	if u, found := c.indexUsername[username]; found {
		return *u
	}

	return User{}
}
