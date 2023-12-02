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
	defer func() {
		if err != nil {
			logrus.Errorln("canal run failed with error, ", err.Error())
		}
	}()
	record, err := dao_pipe.GetRecord(r.Options.PipeName)
	if err != nil {
		return
	}
	if r.pipe.Mysql.Mode == pipeline.MODE_GTID {
		var canGTID mysql.GTIDSet
		if record != nil && record.Pre != nil && record.Pre.GTIDSet != "" {
			if canGTID, err = mysql.ParseGTIDSet(r.pipe.Mysql.Flavor.YaString(), record.Pre.GTIDSet); err != nil {
				return
			}
		} else {
			if canGTID, err = r.storeNewestGTID(); err != nil {
				return
			}
		}
		r.canal.SetEventHandler(&canalHandler{
			ch:   r.OutChan,
			pipe: r.pipe,
		})
		//go r.canal.StartFromGTID(canGTID)
		go func() {
			var startErr error
			defer func() {
				if startErr != nil {
					event.Event(event2.NewErrorPipeline(r.Options.PipeName, "Start mysql replication error: "+startErr.Error()))
				}
			}()
			if startErr = r.canal.StartFromGTID(canGTID); startErr != nil {
				logrus.WithField("mode", "GTID").Errorln(startErr.Error())
				errCode := mysql.ErrorCode(startErr.Error())
				if errCode == 1236 && r.pipe.FixPosNewest {
					if canGTID, startErr = r.storeNewestGTID(); startErr != nil {
						return
					} else {
						event.Event(event2.NewErrorPipeline(r.Options.PipeName, "Start mysql replication could not find first log file name in binary log index file, set current pipeline binlog postion to newest"))
					}
					if startErr = r.canal.StartFromGTID(canGTID); startErr != nil {
						return
					}
				}
			}
		}()
		return
	}

	if r.pipe.Mysql.Mode == pipeline.MODE_POSITION {
		logrus.Debugln("Run pipeline in mode position", r.Options.PipeName)
		var canPos mysql.Position
		if record != nil && record.Pre != nil {
			canPos = mysql.Position{
				Name: record.Pre.BinlogFile,
				Pos:  record.Pre.BinlogPosition,
			}
		} else {
			if canPos, err = r.storeNewestPosition(); err != nil {
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
			var startErr error
			defer func() {
				if startErr != nil {
					event.Event(event2.NewErrorPipeline(r.Options.PipeName, "Start mysql replication error: "+startErr.Error()))
				}
			}()
			if startErr = r.canal.RunFrom(canPos); startErr != nil {
				logrus.WithField("mode", "file index").Errorln(startErr.Error())
				errCode := mysql.ErrorCode(startErr.Error())
				if errCode == 1236 && r.pipe.FixPosNewest {
					if canPos, startErr = r.storeNewestPosition(); startErr != nil {
						return
					} else {
						event.Event(event2.NewErrorPipeline(r.Options.PipeName, "Start mysql replication could not find first log file name in binary log index file, set current pipeline binlog postion to newest"))
					}
					if startErr = r.canal.RunFrom(canPos); startErr != nil {
						return
					}
				}
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

func (r *Input) storeNewestGTID() (gtidSet mysql.GTIDSet, err error) {
	if gtidSet, err = r.canal.GetMasterGTIDSet(); err != nil {
		return
	}
	err = dao_pipe.UpdateRecord(&pipeline.RecordPosition{
		PipelineName: r.Options.PipeName,
		Pre: &pipeline.Position{
			GTIDSet:      gtidSet.String(),
			PipelineName: r.Options.PipeName,
		},
	})
	return
}

func (r *Input) storeNewestPosition() (pos mysql.Position, err error) {
	pos, err = r.canal.GetMasterPos()
	if err != nil {
		logrus.Errorln(err)
		return
	}
	err = dao_pipe.UpdateRecord(&pipeline.RecordPosition{
		PipelineName: r.Options.PipeName,
		Pre: &pipeline.Position{
			BinlogFile:     pos.Name,
			BinlogPosition: pos.Pos,
			PipelineName:   r.Options.PipeName,
		},
	})
	return
}

// Context returns Input's context
func (r *Input) Context() context.Context {
	return r.ctx
}
