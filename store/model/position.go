package model

import "encoding/json"

type Position struct {
	BinlogFile     string `json:"binlog_file"`
	BinlogPosition string `json:"binlog_position"`
	GTIDSet        string `json:"gtid_set"`
	ServerId       string `json:"server_id"`
	ClientId       string `json:"client_id"`
	PipelineID     string `json:"pipeline_id"`
}

func (s *Position) Key() (key string) {
	return "pipeline/" + s.PipelineID + "/position"
}

func (s *Position) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

func (s *Position) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}
