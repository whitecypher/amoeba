package ws

import (
	"sync"
)

var (
	mutex       = sync.Mutex{}
	connections = make([]*Connection, 0)
)

func Broadcast(m Message) {
	mutex.Lock()
	for _, c := range connections {
		c.Send(m)
	}
	mutex.Unlock()
}
