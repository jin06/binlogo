package scheduler

import "testing"

func TestPipelineBind(t *testing.T) {
	pb := EmptyPipelineBind()
	if pb.Key() == "" {
		t.Fail()
	}
	if pb.Val() == "" {
		t.Fail()
	}
	pb2 := EmptyPipelineBind()
	err := pb2.Unmarshal([]byte(pb.Val()))
	if err != nil {
		t.Error(err)
	}
}
