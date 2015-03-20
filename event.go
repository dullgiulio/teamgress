package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type event struct {
	UUID     string    `json:"id"`
	Time     time.Time `json:"time"`
	User     User      `json:"user"`
	EnvName  string    `json:"environment"`
	EnvStage string    `json:"stage"`
	Cmd      string    `json:"command"`
	Data     string    `json:"data"`
}

func newEvent(str string, conf *conf) *event {
	split := strings.SplitN(str, " ", 6)
	e := &event{}

	e.UUID = split[0]

	if unixTime, err := strconv.ParseInt(split[1], 10, 0); err == nil {
		e.Time = time.Unix(int64(unixTime), 0)
	}

	e.User = conf.getByUsername(split[2])
	env := strings.SplitN(split[3], ".", 2)
	e.EnvName = env[0]
	e.EnvStage = env[1]
	e.Cmd = split[4]
	e.Data = split[5]

	return e
}

func (e *event) String() string {
	return fmt.Sprintf("%s: %s.%s '%s': [%s] %s", e.Time.Format(time.RFC3339), e.EnvName, e.EnvStage, &e.User, e.Cmd, e.Data)
}

func (e *event) JSON() string {
	if b, err := json.Marshal(e); err == nil {
		return string(b)
	}

	return ""
}
