package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/whitecypher/amoeba/lib/buffer"
	"github.com/whitecypher/amoeba/lib/filter"
	"github.com/whitecypher/amoeba/lib/ws"
)

type empty struct{}

func Socket() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a new websocket connection
		c := ws.NewConnection(w, r)

		// wrap the websocket connection in buffer the output using the BufferedSender
		b := buffer.NewBuffer(time.Second, c.Sender())
		b.Start()

		// retrieve latency and set buffer interval to latency value
		go func() {
			l := c.Latency()
			fmt.Println("setting buffer latency to ", l.String())
			b.SetInterval(l)
		}()

		// filter out anything outsize of 50,50-150,150
		f := filter.NewViewPort(50, 50, 100, 100, b)

		// update the sender to the filtered and buffered sender
		c.SetSender(f)
	})
}
