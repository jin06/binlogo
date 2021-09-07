package job

import "github.com/jin06/binlogo/server/job/engine"

type Server struct {
	Id string
	Engine *engine.Engine
}

func NewServer(id string) *Server{
	return &Server{}
}


func (s *Server) run() error {
	return nil
}



