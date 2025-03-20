package hotreload

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	addr     string
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWS)
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) Notify(event Event) {
	for client := range s.clients {
		if err := client.WriteJSON(event); err != nil {
			client.Close()
			delete(s.clients, client)
		}
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	s.clients[conn] = true
}
