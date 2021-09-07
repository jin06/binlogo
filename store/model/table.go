package model

import "encoding/json"

type Table struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DatabaseID string `json:"database_id"`
	Charset    string `json:"charset"`
	PrimaryKey string `json:"primary_key"`
}

func (s *Table) Key() (key string) {
	return "table/" + s.ID
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
