package quiz

import "server/internal/server"

type Server struct {
	*server.Server
}

func NewQuizServer(s *server.Server) *Server {
	return &Server{s}
}
