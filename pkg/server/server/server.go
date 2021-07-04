package server

type Server struct {
	Id string
	Engine *Engine
}

func NewServer(id string) *Server {
	return &Server{Id: id}
}

func (s *Server) run() error {
	return nil
}
