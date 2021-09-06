package model

import "encoding/json"

type Pipeline struct {
	ID       string `json:"id"`
	MysqlID  string `json:"mysql_id"`
	Database string `json:"database"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (p *Pipeline) Key() (key string){
	return "pipeline/" + p.ID
}

func (p *Pipeline) Val() (val string) {
	b , _ := json.Marshal(p)
	val = string(b)
	return
}

func (p *Pipeline) Unmarshal(val []byte) (err error) {
	err = json.Unmarshal(val, p)
	return
}


