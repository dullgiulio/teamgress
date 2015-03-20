package libteamgress

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Conf struct {
	Users         []*User `json:"users"`
	indexUsername map[string]*User
}

func NewConf() *Conf {
	return &Conf{
		indexUsername: make(map[string]*User),
		Users:         make([]*User, 0),
	}
}

func (c *Conf) LoadFile(filename string) {
	if err := c._loadJSON(filename); err != nil {
		log.Print(err)
	} else {
		c._indexByUsername()
	}
}

func (c *Conf) _loadJSON(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	file := bufio.NewReader(f)

	if err = json.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	return nil
}
