package input

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/event"
	"github.com/jin06/binlogo/pkg/store/dao/dao_pipe"
	event2 "github.com/jin06/binlogo/pkg/store/model/event"
	"github.com/jin06/binlogo/pkg/store/model/node"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

// Input read binlog event from mysql
type Input struct {
	OutChan      chan *message2.Message
	Options      *Options
	canal        *canal.Canal
	eventHandler *canalHandler
	ctx          context.Context
	pipe         *pipeline.Pipeline
	node         *node.Node
}

// Run Input start working
func (r *Input) Run(ctx context.Context) (err error) {
	myCtx, cancel := context.WithCancel(ctx)
	r.ctx = myCtx
	go func() {
		defer func() {
			if err != nil {
				event.Event(event2.NewErrorPipeline(r.Options.PipeName, err.Error()))
			}
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
				return
			}
		}
	}()
	return
}

// New returns a new Input
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
		Addr:                 addr,
		User:                 pipe.Mysql.User,
		Password:             pipe.Mysql.Password,
		ServerID:             pipe.Mysql.ServerId,
		Flavor:               pipe.Mysql.Flavor.YaString(),
		MaxReconnectAttempts: 3,
	}
	r.canal, err = canal.NewCanal(cfg)
	return
}

func (r *Input) runCanal() (err error) {
	record, err := dao_pipe.GetRecord(r.Options.PipeName)
	if err != nil {
		return
	}
	if r.pipe.Mysql.Mode == pipeline.MODE_GTID {
		var canGTID mysql.GTIDSet
		if record != nil {
			if record.Pre != nil {
				if record.Pre.GTIDSet != "" {
					canGTID, err = mysql.ParseGTIDSet(r.pipe.Mysql.Flavor.YaString(), record.Pre.GTIDSet)
					if err != nil {
						return
					}
				}
			}
		}

		if canGTID == nil {
			canGTID, err = r.canal.GetMasterGTIDSet()
			if err != nil {
				return
			}
			err = dao_pipe.UpdateRecord(&pipeline.RecordPosition{
				PipelineName: r.Options.PipeName,
				Pre: &pipeline.Position{
					GTIDSet:      canGTID.String(),
					PipelineName: r.Options.PipeName,
				},
			})
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
				event.Event(event2.NewErrorPipeline(r.Options.PipeName, "Start mysql replication error: "+startErr.Error()))
			}
		}()
		return
	}

	if r.pipe.Mysql.Mode == pipeline.MODE_POSITION {
		logrus.Debugln("Run pipeline in mode position", r.Options.PipeName)
		var canPos *mysql.Position
		if record != nil {
			if record.Pre != nil {
				canPos = &mysql.Position{
					Name: record.Pre.BinlogFile,
					Pos:  record.Pre.BinlogPosition,
				}
			}
		}
		if canPos == nil {
			canPos = &mysql.Position{}
			*canPos, err = r.canal.GetMasterPos()
			if err != nil {
				logrus.Errorln(err)
				return
			}
			err = dao_pipe.UpdateRecord(&pipeline.RecordPosition{
				PipelineName: r.Options.PipeName,
				Pre: &pipeline.Position{
					BinlogFile:     canPos.Name,
					BinlogPosition: canPos.Pos,
					PipelineName:   r.Options.PipeName,
				},
			})
			if err != nil {
				return
			}
		}
		//logrus.Debugln(pos)
		r.canal.SetEventHandler(&canalHandler{
			ch:   r.OutChan,
			pipe: r.pipe,
		})
		//go r.canal.RunFrom(canPos)
		go func() {
			startErr := r.canal.RunFrom(*canPos)
			if startErr != nil {
				fmt.Println(startErr)
				event.Event(event2.NewErrorPipeline(r.Options.PipeName, "Start mysql replication error: "+startErr.Error()))
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

// Context returns Input's context
func (r *Input) Context() context.Context {
	return r.ctx
}
