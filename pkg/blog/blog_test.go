package blog

import (
	"testing"

	"github.com/jin06/binlogo/v2/configs"
	"github.com/sirupsen/logrus"
)

func TestEnv(t *testing.T) {
	Env(configs.ENV_TEST)
	if logrus.GetLevel() != logrus.DebugLevel {
		t.Fail()
	}
	Env(configs.ENV_DEV)
	if logrus.GetLevel() != logrus.DebugLevel {
		t.Fail()
	}
	Env(configs.ENV_PRO)
	if logrus.GetLevel() != logrus.InfoLevel {
		t.Fail()
	}
}
