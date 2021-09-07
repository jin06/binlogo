package model

import "encoding/json"

type Column struct {
	ID         string   `json:"id"`
	Charset    string   `json:"charset"`
	ColumnType string   `json:"column_type"`
	EnumValues []string `json:"enum_values"`
	PipelineId string   `json:"pipeline_id"`
	DatabaseId string   `json:"database_id"`
	TableId    string   `json:"table_id"`
}

func (s *Column) Key() (key string) {
	return "pipeline/" + s.PipelineId + "/database/" + s.DatabaseId + "/table/" + s.TableId + "/column/" + s.ID
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
