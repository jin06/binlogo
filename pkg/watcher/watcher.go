package watcher

import "github.com/jin06/binlogo/pkg/store/model"

type Watcher interface {
	Watch() chan model.Model
}
