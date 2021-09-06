package store

import (
	"github.com/jin06/binlogo/store/etcd"
	"github.com/jin06/binlogo/store/model"
)

func Create(m model.Model) (bool , error){
	return etcd.E.Create(m)
}
