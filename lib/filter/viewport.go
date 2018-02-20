package filter

import (
	"github.com/whitecypher/amoeba/lib/engine"
	"github.com/whitecypher/amoeba/lib/ws"
)

func NewViewPort(x, y, w, h int, sender ws.Sender) *ViewPort {
	return &ViewPort{
		x:      x,
		y:      y,
		w:      w,
		h:      h,
		sender: sender,
	}
}

// ViewPort filter for filtering of unnecessary messages to connections
type ViewPort struct {
	x      int
	y      int
	w      int
	h      int
	sender ws.Sender
}

func (vp *ViewPort) Sender() ws.Sender {
	return vp.sender
}

// Send implements ws.Sender
func (vp *ViewPort) Send(m ws.Message) {
	// if the event is outside the bounds of the clients view, ignore it.
	if e, ok := m.(engine.Event); ok && (e.X < vp.x || e.Y < vp.y || e.X > vp.x+vp.w || e.Y > vp.y+vp.h) {
		return
	}
	vp.sender.Send(m)
}
