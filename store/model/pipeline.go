package model

import "encoding/json"

type Pipeline struct {
	ID       string `json:"id"`
	MysqlID  string `json:"mysql_id"`
	Database string `json:"database"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s *Pipeline) Key() (key string) {
	return "pipeline/" + s.ID
}

func (s *Pipeline) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Pipeline) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
