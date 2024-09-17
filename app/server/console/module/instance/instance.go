package instance

import (
	"errors"

	"github.com/jin06/binlogo/v2/pkg/store/dao/dao_pipe"
)

func GetInstanceByName(name string) (item *Item, err error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	ins, err := dao_pipe.GetInstance(name)
	if err != nil {
		return
	}
	item = &Item{
		Instance: ins,
	}
	return
}
