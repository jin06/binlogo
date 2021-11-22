package configs

import (
	"github.com/spf13/viper"
	"os"
	"reflect"
	"testing"
)

func TestInitViperFromFile(t *testing.T) {
	file := "./binlogo.yaml"
	InitViperFromFile(file)
	if viper.GetString("env") != "production" {
		t.Fail()
	}
	if viper.GetString("cluster.name") != "cluster1" {
		t.Fail()
	}
	if viper.GetString("node.name") != "node1" {
		t.Fail()
	}
	if viper.GetInt("console.port") != 9999 {
		t.Fail()
	}
	if viper.GetString("console.listen") != "0.0.0.0" {
		t.Fail()
	}
	if !reflect.DeepEqual(viper.GetStringSlice("etcd.endpoints"), []string{"localhost:2379"}){
		t.Fail()
	}
}

func TestInitConfigs(t *testing.T) {
	clusterName := "env_cluster_name"
	err := os.Setenv("CLUSTER_NAME", clusterName)
	InitConfigs()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	hostname, err := os.Hostname()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if viper.GetString("node.name") != hostname {
		t.Fail()
	}

	if viper.GetString("cluster.name") != clusterName {
		t.Log(viper.GetString("cluster.name"))
		t.Error("env get cluster.name fail")
		t.Fail()
	}

}
