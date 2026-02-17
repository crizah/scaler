package auth

import "server/internal/server"

type Server struct {
	*server.Server
}

func NewAuthServer(s *server.Server) *Server {
	return &Server{s}
}
