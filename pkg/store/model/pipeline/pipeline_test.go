package pipeline

import "testing"

func TestPipeline(t *testing.T) {
	p := NewPipeline("go_test_pipe")
	if p.Key() == "" {
		t.Fail()
	}
	if p.Val() == "" {
		t.Fail()
	}
	p2 := &Pipeline{}
	err := p2.Unmarshal([]byte(p.Val()))
	if err != nil {
		t.Error(err)
	}
	if p.Name != p2.Name {
		t.Fail()
	}
}

func TestExpectRunAndWith(t *testing.T) {
	p := NewPipeline("go_test_pipe")
	p.Status = STATUS_RUN
	if p.ExpectRun() == false {
		t.Fail()
	}
	WithPipeStatus(STATUS_STOP)(p)
	if p.ExpectRun() == true {
		t.Fail()
	}
	WithPipeStatus(STATUS_RUN)(p)
	WithPipeDelete(true)(p)
	if p.ExpectRun() == true {
		t.Fail()
	}
	WithPipeSafe(&Pipeline{Mysql: &Mysql{
		Address:  "127.0.0.1",
		Port:     0,
		User:     "",
		Password: "",
		ServerId: 0,
		Flavor:   "",
		Mode:     "",
	},
		AliasName: "alias",
		Filters: []*Filter{
			&Filter{
				Type: "",
				Rule: "",
			},
			&Filter{
				Type: "",
				Rule: "",
			},
		},
		Output: &Output{Sender: &Sender{
			Name:     "",
			Type:     "",
			Kafka:    nil,
			Stdout:   nil,
			Http:     nil,
			RabbitMQ: nil,
			Redis:    nil,
		}},
		Remark: "examples",
	})(p)
	if p.Mysql.Address != "127.0.0.1" {
		t.Fail()
	}
	if len(p.Filters) != 2 {
		t.Fail()
	}
	WithPipeMode(MODE_GTID)(p)
	if p.Mysql.Mode != MODE_GTID {
		t.Fail()
	}
	WithAddFilter(&Filter{})(p)
	if len(p.Filters) != 3 {
		t.Fail()
	}
	WithUpdateFilter(0, &Filter{
		Type: "",
		Rule: "db.tbl",
	})(p)
	if p.Filters[0].Rule != "db.tbl" {
		t.Fail()
	}
}
