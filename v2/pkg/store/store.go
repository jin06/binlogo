package store

import (
	"github.com/jin06/binlogo/v2/pkg/store/model"
)

type Store interface {
	Create(m model.Model) (bool, error)
	Update(m model.Model) (bool, error)
	Delete(m model.Model) (bool, error)
	Get(m model.Model) (bool, error)
}

var def Store

func Init() {

}

// Create  deprecated
func Create(m model.Model) (bool, error) {
	return def.Create(m)
}

// Update  deprecated
func Update(m model.Model) (bool, error) {
	return def.Update(m)
}

// Delete deprecated
func Delete(m model.Model) (bool, error) {
	return def.Delete(m)
}

// Get deprecated
func Get(m model.Model) (bool, error) {
	return def.Get(m)
}
