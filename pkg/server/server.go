package server

import "time"

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	// listen on port 8080
	time.Sleep(10 * time.Second)
}
