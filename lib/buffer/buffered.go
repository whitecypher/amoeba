package buffer

import (
	"time"

	"github.com/whitecypher/amoeba/lib/ws"
)

func NewBuffer(interval time.Duration, sender ws.Sender) *Buffer {
	return &Buffer{
		interval: interval,
		buffer:   make([]ws.Message, 0),
		sender:   sender,
	}
}

// Buffer messages to avoid network overload
type Buffer struct {
	interval  time.Duration
	buffer    []ws.Message
	autoflush bool
	sender    ws.Sender
}

func (bs *Buffer) SetInterval(interval time.Duration) {
	bs.interval = interval
}

func (bs *Buffer) Sender() ws.Sender {
	return bs.sender
}

// Flush the buffer to the sender.Sender
func (bs *Buffer) Flush() {
	if len(bs.buffer) == 0 {
		return
	}
	buf := bs.buffer
	bs.buffer = []ws.Message{}
	bs.sender.Send(ws.Messages(buf))
}

// Start flushing the buffer after each interval
func (bs *Buffer) Start() {
	if bs.interval == time.Duration(0) {
		bs.interval = time.Second
	}
	bs.autoflush = true
	go func() {
		for bs.autoflush {
			bs.Flush()
			time.Sleep(bs.interval)
		}
	}()
}

// Stop flushing the buffer after each interval
func (bs *Buffer) Stop() {
	bs.autoflush = false
}

// Send implements sender.Sender interface
func (bs *Buffer) Send(m ws.Message) {
	bs.buffer = append(bs.buffer, m)
}
