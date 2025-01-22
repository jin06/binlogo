package app

import (
	"bytes"
	"fmt"
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
	blog.SetLevel(configs.Default.LogLevel)
	logrus.Infoln("init configs finish")
	if configs.Default.Profile {
		go func() {
			logrus.Println(http.ListenAndServe(fmt.Sprintf(":%d", configs.Default.ProfilePort), nil))
		}()
	}
	return parseFileToConfig(path)
}

func parseFileToConfig(path string) error {
	var buffer bytes.Buffer
	if _, err := buffer.WriteString(path); err != nil {
		return err
	}
	content, err := os.ReadFile(buffer.String())
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, &configs.Default)
}
