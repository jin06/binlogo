package dao_event

import (
	"github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/jin06/binlogo/pkg/util/key"
	"github.com/spf13/viper"
	"net"
	"testing"
	"time"
)
const (
	nodeName = "test_node_test"
	pipeName = "test_pipeline_test"
)

func init() {
	viper.SetDefault("cluster.name", "test_cluster_test")
	viper.SetDefault("node.name", nodeName)
	viper.SetDefault("etcd.endpoints", []string{"127.0.0.1:12379"})
}

func TestUpdate(t *testing.T) {
	e := event.Event{
		Key: key.GetKey(),
		Type: event.INFO,
		ResourceType: event.PIPELINE,
		ResourceName: pipeName,
		FirstTime: time.Now(),
		LastTime: time.Now(),
		Count: 1,
		MessageType: "",
		Message: "test message",
		NodeName: nodeName,
		NodeIP: net.IP{},
	}
	err := Update(&e)
	if err != nil {
		t.Fail()
	}
}

func TestList(t *testing.T) {
	time.Sleep(time.Second)
	list , err := List("pipeline", pipeName)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if len(list) == 0 {
		t.Fail()
	}

}

func TestDeleteRange(t *testing.T) {
	time.Sleep(time.Second)
	list , err := List("pipeline", pipeName)
	if err != nil {
		t.Fail()
	}
	for _,v := range list {
		from := EventPrefix() + "/" + string(v.ResourceType) + "/" + v.ResourceName + "/" + v.Key
		to := EventPrefix() + "/" + string(v.ResourceType) + "/" + v.ResourceName + "/" + v.Key
		deleted , err2 := DeleteRange(from, to)
		if err2 != nil {
			t.Error(err2)
			t.Fail()
		}
		if deleted == 0{
			t.Fail()
		}
	}
}

