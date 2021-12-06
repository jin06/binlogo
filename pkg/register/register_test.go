package register

import (
	"context"
	"github.com/jin06/binlogo/configs"
	"github.com/jin06/binlogo/pkg/etcdclient"
	"github.com/jin06/binlogo/pkg/util/random"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	configs.InitGoTest()
	r := New(
		WithTTL(1),
		WithKey(etcdclient.Prefix()+"/testregister"+random.String()),
		WithData("1111"),
	)
	ctx, cancel := context.WithCancel(context.Background())
	r.Run(ctx)
	time.Sleep(time.Millisecond * 1200)
	cancel()
	time.Sleep(time.Millisecond * 30)
}
