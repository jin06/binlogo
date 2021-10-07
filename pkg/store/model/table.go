package model

import "encoding/json"

type Table struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Charset      string `json:"charset"`
	PrimaryKey   string `json:"primary_key"`
	PipelineName string `json:"pipeline_name"`
	DatabaseName string `json:"database_name"`
}

func (s *Table) Key() (key string) {
	return "pipeline/" + s.PipelineName + "/database/" + s.DatabaseName + "/table/" + s.Name
}

func (s *Table) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Table) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
