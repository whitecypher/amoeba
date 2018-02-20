package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Ping struct{}

func (p Ping) Bytes() []byte {
	return []byte("ping")
}

type Connection struct {
	ws     *websocket.Conn
	sender Sender
}

func (c *Connection) Sender() Sender {
	return c.sender
}

func (c *Connection) Send(m Message) {
	c.sender.Send(m)
}

func (c *Connection) Latency() time.Duration {
	t := time.Now()
	c.sender.Send(Ping{})
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			continue
		}
		if message == "pong" {
			break
		}
	}
	return time.Now().Sub(t)
}

func (c *Connection) SetSender(s Sender) {
	c.sender = s
}

func NewConnection(w http.ResponseWriter, r *http.Request) *Connection {
	ready := make(chan bool)
	done := make(chan error)

	c := &Connection{}

	go websocket.Handler(func(ws *websocket.Conn) {
		//
		// wrap the websocket in a SimpleSender so it implements the Sender interface
		//
		s := &SimpleSender{
			ws: ws,
		}
		c.ws = ws
		c.sender = s
		close(ready)

		// wait for the connection to be closed by the server because of an error or otherwise
		fmt.Println(<-s.done)
		close(done)
	}).ServeHTTP(w, r)

	// wait for websocket connection to be ready
	<-ready

	// add user to connectedUsers collection so we can start sending intermittent updates
	mutex.Lock()
	connections = append(connections, c)
	mutex.Unlock()

	go func() {
		// wait for websocket to close so we can remove it from the clients list
		err := <-done
		if err != nil {
			log.Println("websocket closed: %v", err)
		}

		// remove uid from connected users
		mutex.Lock()
		for i, c := range connections {
			if c != c {
				continue
			}
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
		mutex.Unlock()
	}()

	return c
}
