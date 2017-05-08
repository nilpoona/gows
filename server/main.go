package main

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

func listenServer(config *Config, router *mux.Router) {
	r := regexp.MustCompile("\\.sock")
	if r.MatchString(config.Server.Host) {

		listener, err := net.Listen("unix", config.Server.Host)
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
		address := config.Server.Host + ":" + config.Server.Port
		addr := flag.String("addr", address, "http service address")
		log.Fatal(http.ListenAndServe(*addr, router))
	}
}

func main() {
	flag.Parse()
	router := mux.NewRouter()
	config := NewConfig()
	hm := newHubManager(config.Hub.Worker)

	hm.runAllHub()

	router.HandleFunc("/room/{id:[1-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomId := vars["id"]
		hub := hm.getHubByRoomid(roomId)
		if hub == nil {
			hub = hm.getHub()
		}
		serveWs(hub, w, r, roomId)
	})

	listenServer(config, router)
}
