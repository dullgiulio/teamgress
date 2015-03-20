package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"sync"
)

type conf struct {
	Users         []*User `json:"users"`
	indexUsername map[string]*User
	mux           *sync.Mutex
}

func newConf() *conf {
	return &conf{
		mux:           &sync.Mutex{},
		indexUsername: make(map[string]*User),
		Users:         make([]*User, 0),
	}
}

func (c *conf) loadFile(filename string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if err := c._loadJSON(filename); err != nil {
		log.Print(err)
	} else {
		c._indexByUsername()
	}
}

func (c *conf) _loadJSON(filename string) error {
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
