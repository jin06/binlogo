package store

import (
	"github.com/jin06/binlogo/store/etcd"
	"github.com/jin06/binlogo/store/model"
)

func Create(m model.Model) (bool, error) {
	return etcd.E.Create(m)
}

func Update(m model.Model) (bool, error) {
	return etcd.E.Update(m)
}

func Delete(m model.Model) (bool, error) {
	return etcd.E.Delete(m)
}

func Get(m model.Model) (bool, error) {
	return etcd.E.Get(m)
}
