package configs

import (
	"github.com/spf13/viper"
	"os"
	"testing"
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
	InitConfigs()
	hostname, err := os.Hostname()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if viper.GetString("node.name") != hostname {
		t.Fail()
	}

	if viper.GetString("cluster.name") != clusterName {
		t.Log(os.Getenv("CLUSTER_NAME"))
		t.Log(viper.GetString("cluster.name"))
		t.Error("env get cluster.name fail")
		t.Fail()
	}

}
