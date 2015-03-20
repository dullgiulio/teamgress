package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type event struct {
	uuid     string
	time     time.Time
	user     string
	envName  string
	envStage string
	cmd      string
	data     string
}

func makeEvent(str string) event {
	split := strings.SplitN(str, " ", 6)
	e := event{}

	e.uuid = split[0]

	if unixTime, err := strconv.ParseInt(split[1], 10, 0); err == nil {
		e.time = time.Unix(int64(unixTime), 0)
	}

	e.user = split[2]
	env := strings.SplitN(split[3], ".", 2)
	e.envName = env[0]
	e.envStage = env[1]
	e.cmd = split[4]
	e.data = split[5]

	return e
}

func (e *event) String() string {
	return fmt.Sprintf("%s: %s.%s %s: [%s] %s", e.time.Format(time.RFC3339), e.envName, e.envStage, e.user, e.cmd, e.data)
}
