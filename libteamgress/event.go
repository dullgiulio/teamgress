package libteamgress

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

type Event struct {
	ID      string    `json:"id"`
	Time    time.Time `json:"time"`
	User    User      `json:"user"`
	Project Project   `json:"project"`
	Type    string    `json:"type"`
	Level   string    `json:"level"` // TODO: Should be enum
	Data    string    `json:"data"`
}

func NewEvent() *Event {
	return &Event{}
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s '%s': %s [%s] %s", e.Time.Format(time.RFC3339), &e.Project, &e.User, e.Level, e.Type, e.Data)
}

func (e *Event) Emit() {
	// TODO: Use json.Compact() and other tricks to make it single line?
	fmt.Printf("%s\n", e.ToJSON())
}

func (e *Event) ToJSON() []byte {
	if b, err := json.Marshal(e); err == nil {
		return b
	}

	return []byte("")
}

// Create an event from JSON byte string
func EventFromJSON(data []byte) (*Event, error) {
	e := NewEvent()

	if err := json.Unmarshal(data, e); err != nil {
		return nil, err
	}

	return e, nil
}

// Read from reader and generate generate events on channel
func ReadJSONEvents(r io.Reader, store *Store) {
	reader := bufio.NewReader(r)

	for {
		text, err := reader.ReadBytes('\n')

		switch err {
		case io.EOF:
			break
		case nil:
			e, err := EventFromJSON(text)

			if err != nil {
				log.Print(err)
			} else {
				store.Add(*e)
			}
		}
	}
}
