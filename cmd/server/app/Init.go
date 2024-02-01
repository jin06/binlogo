package app

import (
	"net/http"

	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/blog"
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

func RunMonitor() {

}
