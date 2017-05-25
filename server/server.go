package server

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/gorilla/mux"
)

func shutdown(listener net.Listener) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
		os.Exit(0)

	}()
}

func NewRouter() *mux.Router {
	return mux.NewRouter()
}

func ListenServer(host, port string, router *mux.Router) {
	r := regexp.MustCompile("\\.sock")
	if r.MatchString(host) {

		listener, err := net.Listen("unix", host)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		defer func() {
			if err := listener.Close(); err != nil {
				log.Printf("error: %v", err)
			}
		}()

		shutdown(listener)
		if err := http.Serve(listener, router); err != nil {
			log.Fatalf("error: %v", err)
		}

	} else {
		address := host + ":" + port
		addr := flag.String("addr", address, "http service address")
		log.Fatal(http.ListenAndServe(*addr, router))
	}
}
