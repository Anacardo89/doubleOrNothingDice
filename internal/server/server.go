package server

import (
	"log"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/api"
	"github.com/Anacardo89/doubleOrNothingDice/internal/websocket"

	"github.com/gorilla/mux"
)

type Server struct {
	router          *mux.Router
	authHandler     *api.AuthHandler
	websocketServer *websocket.Server
}

func NewServer(authHandler *api.AuthHandler) *Server {
	router := mux.NewRouter()
	s := &Server{
		router:          router,
		authHandler:     authHandler,
		websocketServer: websocket.NewServer(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	}).Methods("GET")
	s.router.HandleFunc("/ws", s.websocketServer.UpgradeConnToWS).Methods("GET")
	s.router.HandleFunc("/register", s.authHandler.RegisterHandler).Methods("POST")
	s.router.HandleFunc("/login", s.authHandler.LoginHandler).Methods("POST")
	s.router.HandleFunc("/activate", s.authHandler.ActivateHandler).Methods("GET")
	s.router.HandleFunc("/forgot-password", s.authHandler.ForgotPasswordHandler).Methods("POST")
	s.router.HandleFunc("/recover-password", s.authHandler.RecoverPasswordHandler).Methods("POST")
}

func (s *Server) Run(addr string) {
	log.Printf("HTTP server running at %s\n", addr)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		log.Fatal("HTTP server failed:", err)
	}
}
