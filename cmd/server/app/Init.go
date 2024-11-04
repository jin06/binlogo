package app

import (
	"bytes"
	"net/http"
	"os"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/jin06/binlogo/v2/pkg/blog"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

/*
func SetConfig(path string) *Config {
	if path == "" {
		path = GetCurrentDirectory() + "/configs/data.yml"
	}
	conf := &Config{}
	var buffer bytes.Buffer
	buffer.WriteString(path)
	fp := buffer.String()
	if ymlFile, err := ioutil.ReadFile(fp); err != nil {
		panic(err)
	} else if err = yaml.Unmarshal(ymlFile, &conf); err != nil {
		panic(err)
	}
	return conf
}
*/
// Init init run environment
func Init(path string) error {
	if err := configs.Init(path); err != nil {
		return err
	}
	blog.Env(configs.ENV)
	logrus.Infoln("init configs finish")
	if configs.ENV == configs.ENV_DEV {
		go func() {
			logrus.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	err := parseFileToConfig(path)
	return err
}

func parseFileToConfig(path string) error {
	var buffer bytes.Buffer
	var err error
	var content []byte
	buffer.WriteString(path)
	fp := buffer.String()
	if content, err = os.ReadFile(fp); err != nil {
		return err
	}
	conf := configs.Config{}
	if err = yaml.Unmarshal(content, &conf); err != nil {
		return err
	}
	configs.Default = conf
	return nil
}
