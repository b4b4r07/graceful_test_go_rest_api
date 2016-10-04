package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var now = time.Now()

func main() {
	log.Printf("start pid %d\n", os.Getpid())
	s := &http.Server{Addr: ":8080", Handler: newHandler()}
	s.ListenAndServe()
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
