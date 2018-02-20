package ws

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Sender interface {
	Sender() Sender
	Send(Message)
}

type SimpleSender struct {
	done chan error
	ws   *websocket.Conn
}

func (s *SimpleSender) Sender() Sender {
	return s
}

func (s *SimpleSender) Send(m Message) {
	_, err := s.ws.Write(m.Bytes())
	if err != nil {
		fmt.Println("failed to send message: ", err)
		s.done <- fmt.Errorf("failed to send message %v", err)
	}
}
