package store

import (
	"context"

	"github.com/jin06/binlogo/v2/pkg/store/etcd"
	"github.com/jin06/binlogo/v2/pkg/store/model"
	"github.com/spf13/viper"
)

var def Store

type Store interface {
	Create(ctx context.Context, m model.Model) (bool, error)
	Update(ctx context.Context, m model.Model) (bool, error)
	Delete(ctx context.Context, m model.Model) (bool, error)
	Get(ctx context.Context, m model.Model) (bool, error)
}

func Init() {
	switch viper.GetString("store.type") {
	case "etcd":
		{

		}
	case "redis":
		{
		}
	default:
		{
			panic("Unsupport store type")
		}
	}
}

// Create  deprecated
func Create(m model.Model) (bool, error) {
	return etcd.E.Create(m)
}

// Update  deprecated
func Update(m model.Model) (bool, error) {
	return etcd.E.Update(m)
}

// Delete deprecated
func Delete(m model.Model) (bool, error) {
	return etcd.E.Delete(m)
}

// Get deprecated
func Get(m model.Model) (bool, error) {
	return etcd.E.Get(m)
}
