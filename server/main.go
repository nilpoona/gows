package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8088", "http service address")

func shutdown(listener net.Listener) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
		os.Exit(1)

	}()
}

func main() {
	flag.Parse()
	hub := newHub()
	r := mux.NewRouter()
	go hub.run()
	r.HandleFunc("/room/{id:[1-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serveWs(hub, w, r, vars["id"])
	})

	listener, err := net.Listen("unix", "/tmp/ws.sock")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("error: %v", err)
		}
	}()

	if err := http.Serve(listener, r); err != nil {
		log.Fatalf("error: %v", err)
	}
	/*
		err := http.ListenAndServe(*addr, r)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	*/
}
