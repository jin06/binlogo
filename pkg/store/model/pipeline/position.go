package pipeline

import "encoding/json"

type Position struct {
	BinlogFile     string `json:"binlog_file"`
	BinlogPosition uint32 `json:"binlog_position"`
	GTIDSet        string `json:"gtid_set"`
	PipelineName   string `json:"pipeline_name"`
}

func (s *Position) Key() (key string) {
	return "run/position/" + s.PipelineName
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
