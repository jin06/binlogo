package model

import "encoding/json"

type Mysql struct {
	Address    string `json:"address"`
	Port       int `json:"post"`
	User       string `json:"user"`
	Password   string `json:"password"`
	PipelineId string `json:"pipeline_id"`
	ServerId   int `json:"server_id"`
}

func (s *Mysql) Key() (key string) {
	return "pipeline/" + s.PipelineId + "/mysql"
}

func (s *Mysql) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Mysql) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
