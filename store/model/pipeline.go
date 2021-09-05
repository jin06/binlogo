package model

type Pipeline struct {
	ID       string `json:"id"`
	MysqlID  string
	Database string
	Name     string
	Password string
}

func (p *Pipeline) Key() (key string){
	return
}

func (p *Pipeline) Val() (val string) {
	return
}


