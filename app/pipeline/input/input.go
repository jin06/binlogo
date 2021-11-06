package input

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/blog"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Input struct {
	OutChan      chan *message2.Message
	Options      *Options
	canal        *canal.Canal
	eventHandler *canalHandler
}

func (r *Input) Run(ctx context.Context) (err error) {
	err = r.runCanal()
	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					r.canal.Close()
					return
				}
			case <-time.Tick(5):
				{
					blog.Debugln("input tick")
				}
			default:
			}
		}
	}()
	return
}

func New(opts ...Option) (input *Input, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	input = &Input{
		Options: options,
	}
	err = input.init()

	return
}

func (r *Input) init() (err error) {
	addr := fmt.Sprintf("%s:%s", r.Options.Mysql.Address, strconv.Itoa(int(r.Options.Mysql.Port)))
	cfg := &canal.Config{
		Addr:     addr,
		User:     r.Options.Mysql.User,
		Password: r.Options.Mysql.Password,
		ServerID: r.Options.Mysql.ServerId,
		Flavor:   r.Options.Mysql.Flavor,
	}
	r.canal, err = canal.NewCanal(cfg)
	if err != nil {
		return
	}
	return
}

func (r *Input) runCanal() (err error) {
	go func() {
		if r.Options.Pipeline.Mode == pipeline.MODE_GTID {
			return
		}

		if r.Options.Pipeline.Mode == pipeline.MODE_POSTION {
			blog.Debugln("Run pipeline in mode position", r.Options.Pipeline.Name)
			pos := mysql.Position{}
			binlogFile := r.Options.Position.BinlogFile
			if binlogFile == "" {
				logrus.Warn("Empty binlog file")
				pos, err = r.canal.GetMasterPos()
				return
			} else {
				binlogPos := r.Options.Position.BinlogPosition
				pos = mysql.Position{
					binlogFile,
					binlogPos,
				}
			}
			r.canal.SetEventHandler(&canalHandler{
				ch:           r.OutChan,
				positionFile: pos.Name,
				pipe:         r.Options.Pipeline,
			})
			err = r.canal.RunFrom(pos)
			return
		}

		err = errors.New("invalid mode")
		return
	}()
	time.Sleep(1)
	return
}
