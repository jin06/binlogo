package dao_pipe

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/store/model/pipeline"
	"github.com/jin06/binlogo/v2/pkg/util/random"
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
