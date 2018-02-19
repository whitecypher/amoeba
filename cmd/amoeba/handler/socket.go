package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Message interface {
	Bytes() []byte
}

type Messages []Message

func (m Messages) Bytes() []byte {
	r := [][]byte{}
	for _, i := range m {
		r = append(r, i.Bytes())
	}
	return bytes.Join(r, []byte("\n"))

}

type Event struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Direction int `json:"direction"`
	Velocity  int `json:"velocity"`
	Decay     int `json:"decay"`
}

func (e Event) Bytes() []byte {
	r, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
	}
	return r
}

type Sender interface {
	Send(Message)
}

type Client struct {
	ws Sender
}

func (c *Client) Send(m Message) {
	c.ws.Send(m)
}

type RawSocket struct {
	done chan error
	ws   *websocket.Conn
}

func (rs *RawSocket) Send(m Message) {
	_, err := rs.ws.Write(m.Bytes())
	if err != nil {
		fmt.Println("failed to send message: ", err)
		rs.done <- fmt.Errorf("failed to send message %v", err)
	}
}

type FilteredSender struct {
	x  int
	y  int
	w  int
	h  int
	ws Sender
}

func (fs *FilteredSender) Send(m Message) {
	// if the event is outside the bounds of the clients view, ignore it.
	if e, ok := m.(Event); ok && (e.X < fs.x || e.Y < fs.y || e.X > fs.x+fs.w || e.Y > fs.y+fs.h) {
		return
	}
	fs.ws.Send(m)
}

type BufferedSender struct {
	interval time.Duration
	buffer   []Message
	enabled  bool
	ws       Sender
}

func (bs *BufferedSender) Start() {
	if bs.interval == time.Duration(0) {
		bs.interval = time.Second
	}
	bs.enabled = true
	go func() {
		for bs.enabled {
			if len(bs.buffer) == 0 {
				continue
			}
			buf := bs.buffer
			bs.buffer = []Message{}
			bs.ws.Send(Messages(buf))
			time.Sleep(bs.interval)
		}
	}()
}

func (bs *BufferedSender) Stop() {
	bs.enabled = false
}

func (bs *BufferedSender) Send(m Message) {
	bs.buffer = append(bs.buffer, m)
}

type empty struct{}

var (
	mutex   = sync.Mutex{}
	clients = make([]*Client, 0)
)

func Socket() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ready := make(chan bool)
		done := make(chan error)

		client := &Client{}

		//t := time.Now()

		go websocket.Handler(func(ws *websocket.Conn) {
			//
			// wrap the websocket in a RawSocket so it implements the Sender interface
			//
			rs := &RawSocket{
				ws: ws,
			}
			client.ws = rs
			close(ready)

			// wait for the connection to be closed by the server because of an error or otherwise
			fmt.Println(<-rs.done)
			close(done)
		}).ServeHTTP(w, r)

		// wait for websocket connection to be ready
		<-ready

		// wrap the websocket connection in buffer the output using the BufferedSender
		b := &BufferedSender{
			// set the interval the time it takes to activate the websocket
			interval: time.Second, // time.Now().Sub(t)
			ws:       client.ws,
		}
		// filter out anything outsize of 50,50-150,150
		client.ws = &FilteredSender{
			x:  50,
			y:  50,
			w:  100,
			h:  100,
			ws: b,
		}
		b.Start()

		// add user to connectedUsers collection so we can start sending intermittent updates
		mutex.Lock()
		clients = append(clients, client)
		mutex.Unlock()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

		// wait for websocket to close so we can remove it from the clients list
		err := <-done
		if err != nil {
			log.Println("websocket closed: %v", err)
		}

		// remove uid from connected users
		mutex.Lock()
		for i, c := range clients {
			if c != client {
				continue
			}
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
		mutex.Unlock()

	})
}

func Broadcast(e Event) {
	mutex.Lock()
	for _, c := range clients {
		c.Send(e)
	}
	mutex.Unlock()
}
