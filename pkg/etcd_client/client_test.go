package etcd_client

import (
	"context"
	"github.com/spf13/viper"
	"testing"
)

func TestNew(t *testing.T) {
	endPoints := []string{"127.0.0.1:12379"}
	viper.SetDefault("etcd.endpoints", endPoints)
	viper.SetDefault("etcd.password", "")
	viper.GetStringSlice("etcd.endpoints")
	cli, err := New()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	res, err := cli.Status(context.Background(), endPoints[0])
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf(res.Header.String())

}
