package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alphacentaurigames/conquest-alpha-go/lib/routing"
	"github.com/whitecypher/amoeba/cmd/amoeba/handler"
	"github.com/whitecypher/logr"
)

func main() {
	r := routing.New()
	r.Middleware()

	r.HandleFunc("/", handler.File("public/index.html"))
	r.HandleFunc("/main.js", handler.File("public/main.js"))
	r.HandleFunc("/main.css", handler.File("public/main.css"))
	r.HandleFunc("/ws", handler.Socket())

	// create the server
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	// listen for SIGKILL
	sig := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		// start server
		errChan <- server.ListenAndServe()
	}()

	go func() {
		// send messages
		for {
			e := handler.Event{
				X:         rand.Intn(200),
				Y:         rand.Intn(200),
				Direction: rand.Intn(360),
				Velocity:  rand.Intn(10),
				Decay:     1,
			}
			handler.Broadcast(e)
			time.Sleep(time.Second / time.Duration(rand.Intn(9)+1))
		}
	}()

	var err error
	select {
	case <-sig:
		// shutdown the server
		os.Stdout.WriteString("\n")
		logr.Info("shutting down the web server...")
		timeout := time.Second * 10
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		err = server.Shutdown(ctx)
		if err == context.DeadlineExceeded {
			logr.Infof("web server stopped forcefully after %s", timeout)
		} else if err != nil {
			logr.Infof("web server stopped with err: %s", err)
		} else {
			logr.Info("web server stopped gracefully")
		}
	case err = <-errChan:
		// nothing to do
	}

	fmt.Println(err)
}
