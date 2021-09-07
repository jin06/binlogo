package model

import "encoding/json"

type Database struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	SchemasID  string `json:"schema_id"`
	Charset    string `json:"charset"`
	PipelineId string `json:"pipeline_id"`
}

func (s *Database) Key() string {
	return "pipeline/" + s.PipelineId + "/database/" + s.ID
}

func (s *Database) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Database) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
