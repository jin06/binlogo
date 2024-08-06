package configs

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestInitViperFromFile(t *testing.T) {
	file := "./binlogo.yaml"
	Init(file)
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
	if !reflect.DeepEqual(viper.GetStringSlice("etcd.endpoints"), []string{"localhost:2379"}) {
		t.Fail()
	}
}
