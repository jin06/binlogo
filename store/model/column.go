package model

import "encoding/json"

type Column struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Charset      string   `json:"charset"`
	Type         string   `json:"type"`
	EnumValues   []string `json:"enum_values"`
	PipelineName string   `json:"pipeline_name"`
	DatabaseName string   `json:"database_name"`
	TableName    string   `json:"table_name"`
	Signed       string   `json:"signed"`
}

func (s *Column) Key() (key string) {
	return "pipeline/" + s.PipelineName + "/database/" + s.DatabaseName + "/table/" + s.TableName + "/column/" + s.Name
}

func (s *Column) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Column) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
