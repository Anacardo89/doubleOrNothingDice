package server

import (
	"log"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/websocket"
)

type Server struct {
	websocketServer *websocket.Server
}

func NewServer() *Server {
	wsServer := websocket.NewServer()
	return &Server{
		websocketServer: wsServer,
	}
}

func (s *Server) Run(addr string) {
	http.HandleFunc("/ws", s.websocketServer.UpgradeConnToWS)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("WebSocket server is running"))
	})

	log.Printf("HTTP server running at %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("HTTP server failed:", err)
	}
}
