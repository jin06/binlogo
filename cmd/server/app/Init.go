package app

import (
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Init init run environment
func Init(file string) {
	configs.InitConfigs()
	configs.InitViperFromFile(file)
	//etcd2.DefaultETCD()
	blog.Env(configs.ENV)
	logrus.Infoln("init configs finish")
	if configs.ENV == configs.ENV_DEV {
		go func() {
			logrus.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

}
