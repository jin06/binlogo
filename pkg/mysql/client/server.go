package client

import (
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
)

func newConn(mysql *pipeline.Mysql) (conn *client.Conn, err error) {
	conn, err = client.Connect(mysql.Address, mysql.User, mysql.Password, "")
	if err != nil {
		return
	}
	return
}

func MasterStatus(mysql *pipeline.Mysql) (err error){
	conn, err := newConn(mysql)
	if err != nil {
		return
	}
	defer conn.Close()
	result, err := conn.Execute("show master status")
	if err != nil {
		return
	}
	_  = result
	return
}
