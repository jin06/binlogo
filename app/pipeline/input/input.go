package input

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
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
	ctx          context.Context
}

func (r *Input) Run(ctx context.Context) (err error) {
	myCtx, cancel := context.WithCancel(ctx)
	r.ctx = myCtx
	go func() {
		defer func() {
			if r.canal != nil {
				r.canal.Close()
			}
			cancel()
		}()
		for {
			err = r.prepareCanal()
			if err != nil {
				logrus.Errorln(err)
				continue
			}
			go r.runCanal()
			select {
			case <-ctx.Done():
				{
					return
				}
			case <-r.canal.Ctx().Done():
				{
					break
				}
			}
			time.Sleep(time.Second)
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
	//err = input.init()

	return
}

func (r *Input) prepareCanal() (err error) {
	pipe, err := dao_pipe.GetPipeline(r.Options.Pipeline.Name)
	if err != nil {
		return
	}
	addr := fmt.Sprintf("%s:%s", pipe.Mysql.Address, strconv.Itoa(int(pipe.Mysql.Port)))
	cfg := &canal.Config{
		Addr:     addr,
		User:     pipe.Mysql.User,
		Password: pipe.Mysql.Password,
		ServerID: pipe.Mysql.ServerId,
		Flavor:   pipe.Mysql.Flavor,
	}
	r.canal, err = canal.NewCanal(cfg)
	return
}

func (r *Input) runCanal() (err error) {
	if r.Options.Pipeline.Mysql.Mode == pipeline.MODE_GTID {
		return
	}

	if r.Options.Pipeline.Mysql.Mode == pipeline.MODE_POSTION {
		logrus.Debugln("Run pipeline in mode position", r.Options.Pipeline.Name)
		pos := mysql.Position{}
		binlogFile := r.Options.Position.BinlogFile
		if binlogFile == "" {
			logrus.Warn("Empty binlog file")
			pos, err = r.canal.GetMasterPos()
			if err != nil {
				logrus.Errorln(err)
				return
			}
		} else {
			binlogPos := r.Options.Position.BinlogPosition
			pos = mysql.Position{
				binlogFile,
				binlogPos,
			}
		}
		//logrus.Debugln(pos)
		r.canal.SetEventHandler(&canalHandler{
			ch:           r.OutChan,
			positionFile: pos.Name,
			pipe:         r.Options.Pipeline,
		})
		err = r.canal.RunFrom(pos)
		return
	}

	err = errors.New("invalid mode")
	//if err != nil {
	//	logrus.Errorln("Run canal error: ", err)
	//}
	return
}

func (r *Input) Context() context.Context {
	return r.ctx
}
