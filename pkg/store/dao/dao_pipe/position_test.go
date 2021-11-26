package dao_pipe

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
)

func TestPosition(t *testing.T) {
	configs.InitGoTest()
	pipeName := "go_test_pipeline" + random.String()
	p := &pipeline.Position{
		BinlogFile:     "",
		BinlogPosition: 0,
		GTIDSet:        "",
		PipelineName:   pipeName,
	}
	err := UpdatePosition(p)
	if err != nil {
		t.Error(err)
	}
	_, err = UpdatePositionSafe(pipeName,
		pipeline.WithBinlogFile(""),
		pipeline.WithPos(uint32(0)),
		pipeline.WithGTIDSet(""),
	)
	if err != nil {
		t.Error(err)
	}
	p2, err := GetPosition(pipeName)
	if err != nil {
		t.Error(err)
	}
	if p.PipelineName != p2.PipelineName {
		t.Error("pipeline name different")
	}
	ok, err := DeletePosition(pipeName)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("delete not ok ")
	}
}
