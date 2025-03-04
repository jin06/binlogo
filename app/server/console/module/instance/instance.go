package instance

import (
	"context"
	"errors"

	"github.com/jin06/binlogo/v2/pkg/store/dao"
)

func GetInstanceByName(name string) (item *Item, err error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	ins, err := dao.GetInstance(context.Background(), name)
	if err != nil {
		return
	}
	item = &Item{
		Instance: ins,
	}
	return
}
