package model

import "encoding/json"

type Server struct {
	ID         string `json:"id"`
	Address    string `json:"address"`
	Port       string `json:"post"`
	User       string `json:"user"`
	Password   string `json:"password"`
	PipelineId string `json:"pipeline_id"`
}

func (s *Server) Key() (key string) {
	return "pipeline/" + s.PipelineId + "/server"
}

func (s *Server) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Server) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
