package store

import (
	etcd2 "github.com/jin06/binlogo/pkg/store/etcd"
	model2 "github.com/jin06/binlogo/pkg/store/model"
)

func Create(m model2.Model) (bool, error) {
	return etcd2.E.Create(m)
}

func Update(m model2.Model) (bool, error) {
	return etcd2.E.Update(m)
}

func Delete(m model2.Model) (bool, error) {
	return etcd2.E.Delete(m)
}

func Get(m model2.Model) (bool, error) {
	return etcd2.E.Get(m)
}
