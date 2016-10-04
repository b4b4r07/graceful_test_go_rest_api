package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/braintree/manners"
	"github.com/lestrrat/go-server-starter-listener"
)

var now = time.Now()

func main() {
	log.Printf("start pid %d\n", os.Getpid())

	signal_chan := make(chan os.Signal)
	signal.Notify(signal_chan, syscall.SIGTERM)
	go func() {
		for {
			s := <-signal_chan
			if s == syscall.SIGTERM {
				manners.Close()
			}
		}
	}()

	l, err := ss.NewListener()
	if l == nil || err != nil {
		// Fallback if not running under Server::Starter
		l, err = net.Listen("tcp", ":8080")
		if err != nil {
			panic("Failed to listen to port 8080")
		}
	}

	manners.Serve(l, newHandler())
}

// https://github.com/facebookgo/grace/blob/master/gracedemo/demo.go から一部拝借
func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/sleep/", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		time.Sleep(duration)
		fmt.Fprintf(
			w,
			"started at %s slept for %d nanoseconds from pid %d.\n",
			now,
			duration.Nanoseconds(),
			os.Getpid(),
		)
	})
	return mux
}
