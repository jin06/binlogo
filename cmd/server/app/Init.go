package app

import (
	"net/http"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/blog"
	"github.com/sirupsen/logrus"
)

// Init init run environment
func Init(file string) {
	configs.Init(file)
	blog.Env(configs.ENV)
	logrus.Infoln("init configs finish")
	if configs.ENV == configs.ENV_DEV {
		go func() {
			logrus.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}
