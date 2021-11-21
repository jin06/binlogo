package input

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/event"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Input struct {
	OutChan      chan *message2.Message
	Options      *Options
	canal        *canal.Canal
	eventHandler *canalHandler
	ctx          context.Context
	pipe         *pipeline.Pipeline
	node         *node.Node
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
		if r.canal == nil {
			err = r.prepareCanal()
			if err != nil {
				return
			}
			err = r.runCanal()
			if err != nil {
				return
			}
		}
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
	return
}

func (r *Input) prepareCanal() (err error) {
	pipe, err := dao_pipe.GetPipeline(r.Options.PipeName)
	if err != nil {
		return
	}
	if pipe == nil {
		err = errors.New("pipeline not found")
		return
	}
	r.pipe = pipe

	addr := fmt.Sprintf("%s:%s", pipe.Mysql.Address, strconv.Itoa(int(pipe.Mysql.Port)))
	cfg := &canal.Config{
		Addr:     addr,
		User:     pipe.Mysql.User,
		Password: pipe.Mysql.Password,
		ServerID: pipe.Mysql.ServerId,
		Flavor:   pipe.Mysql.Flavor.YaString(),
	}
	r.canal, err = canal.NewCanal(cfg)
	return
}

func (r *Input) runCanal() (err error) {
	pos, err := dao_pipe.GetPosition(r.Options.PipeName)
	if err != nil {
		return
	}

	if r.pipe.Mysql.Mode == pipeline.MODE_GTID {
		var canGTID mysql.GTIDSet
		if pos != nil {
			if pos.GTIDSet != "" {
				canGTID, err = mysql.ParseGTIDSet(r.pipe.Mysql.Flavor.YaString(), pos.GTIDSet)
				if err != nil {
					return
				}
			}
		}

		if canGTID == nil {
			pos = &pipeline.Position{}
			canGTID, err = r.canal.GetMasterGTIDSet()
			if err != nil {
				return
			}
		}
		r.canal.SetEventHandler(&canalHandler{
			ch:   r.OutChan,
			pipe: r.pipe,
		})
		//go r.canal.StartFromGTID(canGTID)
		go func() {
			startErr := r.canal.StartFromGTID(canGTID)
			if startErr != nil {
				fmt.Println("123",  startErr)
				event.Event(event2.NewErrorPipeline(r.pipe.Name, "Start mysql replication error: " + startErr.Error()))
			}
		}()
		return
	}

	if r.pipe.Mysql.Mode == pipeline.MODE_POSTION {
		logrus.Debugln("Run pipeline in mode position", r.Options.PipeName)
		var canPos mysql.Position
		if pos == nil {
			pos = &pipeline.Position{}
			logrus.Warn("Empty binlog file")
			canPos, err = r.canal.GetMasterPos()
			if err != nil {
				logrus.Errorln(err)
				return
			}
		} else {
			canPos = mysql.Position{
				pos.BinlogFile,
				pos.BinlogPosition,
			}
		}
		//logrus.Debugln(pos)
		r.canal.SetEventHandler(&canalHandler{
			ch:   r.OutChan,
			pipe: r.pipe,
		})
		//go r.canal.RunFrom(canPos)
		go func() {
			startErr := r.canal.RunFrom(canPos)
			if startErr != nil {
				fmt.Println(startErr)
				event.Event(event2.NewErrorPipeline(r.pipe.Name, "Start mysql replication error: " + startErr.Error()))
			}
		}()
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
