package pipeline

import "testing"

func TestPosition(t *testing.T) {
	p := &Position{
		BinlogFile:     "mysql-bin.000004",
		BinlogPosition: 17561,
		GTIDSet:        "045c649a-408d-11ec-ae21-0242ac110006:1-64",
		PipelineName:   "go_test_pipeline",
	}
	if p.Key() == "" {
		t.Fail()
	}
	p2 := &Position{}
	err := p2.Unmarshal([]byte(p.Val()))
	if err != nil {
		t.Error(err)
	}
	if p.PipelineName != p2.PipelineName {
		t.Fail()
	}
	file := "mysql-bin.000005"
	pos := uint32(1012)
	gtid := "045c649a-408d-11ec-ae21-0242ac110006:1-200"
	WithBinlogFile(file)(p)
	if p.BinlogFile != file {
		t.Fail()
	}
	WithPos(pos)(p)
	if p.BinlogPosition != pos {
		t.Fail()
	}
	WithGTIDSet(gtid)(p)
	if p.GTIDSet != gtid {
		t.Fail()
	}
}
