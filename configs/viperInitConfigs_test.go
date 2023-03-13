package configs

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfigs(t *testing.T) {
	clusterName := "env_cluster_name"
	err := os.Setenv("CLUSTER_NAME", clusterName)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	viper.Set("node.name", nil)
	viper.Set("cluster.name", nil)
	InitGoTest()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(viper.GetString("node.name"))
	t.Log(viper.GetString("cluster.name"))
}
