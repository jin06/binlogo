package model

import "encoding/json"

type Pipeline struct {
	Name      string   `json:"name"`
	AliasName string   `json:"aliasName"`
	Mysql     *Mysql    `json:"mysql"`
	Filters   []*Filter `json:"filters"`
	Output    *Output   `json:"output"`
}

func (s *Pipeline) Key() (key string) {
	return "pipeline/" + s.Name
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
