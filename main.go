package main

import (
	"flag"
	"net/http"

	"github.com/YutakaHorikawa/gows/config"
	"github.com/YutakaHorikawa/gows/hub"
	"github.com/YutakaHorikawa/gows/server"
	"github.com/YutakaHorikawa/gows/ws"

	"github.com/gorilla/mux"
)

func main() {
	flag.Parse()
	router := server.NewRouter()
	config := config.NewConfig()
	hm := hub.NewHubManager(config.Hub.Worker)

	hm.RunAllHub()

	router.HandleFunc("/room/{id:[1-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomId := vars["id"]
		hub := hm.GetHubByRoomid(roomId)
		if hub == nil {
			hub = hm.GetHub()
		}
		ws.ServeWs(hub, w, r, roomId)
	})

	server.ListenServer(config.Server.Host, config.Server.Port, router)
}
