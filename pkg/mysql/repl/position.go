package repl

import "github.com/go-mysql-org/go-mysql/mysql"

type Position struct {
	File string
	Position uint32
}

func (p *Position) BinlogPosition() mysql.Position{
	return mysql.Position{
		Name: p.File, Pos: p.Position,
	}
}
