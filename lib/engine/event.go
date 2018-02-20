package engine

import (
	"encoding/json"
	"log"
)

// Item event
type Event struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Direction int `json:"direction"`
	Velocity  int `json:"velocity"`
	Decay     int `json:"decay"`
}

// Bytes implements ws.Message
func (e Event) Bytes() []byte {
	r, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return r
}
