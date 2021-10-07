package input

import (
	"context"
	"fmt"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/store/model"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/sirupsen/logrus"
)

type Input struct {
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer
	OutChan  chan *message2.Message
	Options  *Options
}

func (r *Input) Run() (err error) {
	if err = r.connect(); err != nil {
		return
	}
	err = r.handle()
	return
}

func (r *Input) connect() (err error) {
	binlogFile := r.Options.Position.BinlogFile
	if binlogFile == "" {
		logrus.Warn("Empty binlog file")
		//return errors.New("empty binlog file")
	}
	binlogPos := r.Options.Position.BinlogPosition

	pos := mysql.Position{
		binlogFile,
		binlogPos,
	}
	streamer, err := r.syncer.StartSync(pos)
	if err != nil {
		return
	}
	r.streamer = streamer
	return
}

func (r *Input) doHandle() {
	for {
		logrus.Debug("Sync binlog from mysql")
		ctx := context.Background()
		e, er := r.streamer.GetEvent(ctx)
		if er != nil {
			logrus.Error(er)
			continue
		}
		logrus.Debug("Binlog event header : ", e.Header)
		msg, err := inputMessage(e)
		pos := r.syncer.GetNextPosition()
		fmt.Println(pos)
		msg.BinlogPosition = &model.Position{
			BinlogFile:     pos.Name,
			BinlogPosition: pos.Pos,
			GTIDSet:        r.Options.Position.GTIDSet,
			PipelineName:   r.Options.Position.PipelineName,
		}
		if err != nil {
			panic(err)
		}
		if msg != nil {
			r.OutChan <- msg
		} else {
			logrus.Debug("The event is not a data change event")
		}
	}
	return
}

func (r *Input) handle() (err error) {
	go r.doHandle()
	return
}

func (r *Input) DataLine() chan *message2.Message {
	return r.OutChan
}

func New(opts ...Option) (input *Input, err error) {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	input = &Input{
		Options: options,
	}
	err = input.Init()

	return
}

func (r *Input) Init() (err error) {
	cfg := replication.BinlogSyncerConfig{
		ServerID: r.Options.Mysql.ServerId,
		Flavor:   r.Options.Mysql.Flavor,
		Host:     r.Options.Mysql.Address,
		Port:     r.Options.Mysql.Port,
		User:     r.Options.Mysql.User,
		Password: r.Options.Mysql.Password,
	}

	r.syncer = replication.NewBinlogSyncer(cfg)
	//r.Ch = make(chan *message.Message, 100000)
	return
}

func (r *Input) newestPosition() {
}
