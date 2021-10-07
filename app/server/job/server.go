package job

import (
	engine2 "github.com/jin06/binlogo/app/server/job/engine"
)

type Server struct {
	Id string
	Engine *engine2.Engine
}

func NewServer(id string) *Server {
	return &Server{}
}


func (s *Server) run() error {
	return nil
}



