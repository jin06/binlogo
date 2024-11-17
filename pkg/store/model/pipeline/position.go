package pipeline

import (
	"encoding/json"
)

// Position mysql replication position
type Position struct {
	BinlogFile     string `json:"binlog_file" redis:"binlog_file"`
	BinlogPosition uint32 `json:"binlog_position" redis:"binlog_position"`
	GTIDSet        string `json:"gtid_set" redis:"gtid_set"`
	PipelineName   string `json:"pipeline_name" redis:"pipeline_name"`
	TotalRows      int    `json:"total_rows" redis:"total_rows"`
	ConsumeRows    int    `json:"consume_rows" redis:"consume_rows"`
}

// Key get etcd key prefix
func (s *Position) Key() (key string) {
	return "run/position/" + s.PipelineName
}

// Val get position json data
func (s *Position) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}

// Unmarshal unmarshal json data to object
func (s *Position) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, s)
	return
}

// OptionPosition Position options
type OptionPosition func(position *Position)

// WithBinlogFile sets binlog file to OptionPosition
func WithBinlogFile(b string) OptionPosition {
	return func(position *Position) {
		position.BinlogFile = b
	}
}

// WithPos sets binlog position to OptionPosition
func WithPos(p uint32) OptionPosition {
	return func(position *Position) {
		position.BinlogPosition = p
	}
}

// WithGTIDSet sets GTIDSet to OptionPosition
func WithGTIDSet(g string) OptionPosition {
	return func(position *Position) {
		position.GTIDSet = g
	}
}

// Reset reset position
func (p *Position) Reset() {
	p.ConsumeRows = 0
	p.TotalRows = 0
	p.BinlogFile = ""
	p.GTIDSet = ""
	p.PipelineName = ""
	p.BinlogPosition = 0
}
