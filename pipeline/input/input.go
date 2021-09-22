package input

import (
	"context"
	"errors"
	"fmt"
	"github.com/jin06/binlogo/pipeline/message"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/sirupsen/logrus"
)

type Input struct {
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer
	Ch       chan *message.Message
	Options  *Options
}

func (r *Input) Start() (err error) {
	err = r.connect()
	if err != nil {
		return
	}
	err = r.handle()
	return
}

func (r *Input) connect() (err error) {
	binlogFile := r.Options.Position.BinlogFile
	if binlogFile == "" {
		return errors.New("empty binlog file")
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

func (r *Input) sync() (err error) {
	go func() {
		//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		for {
			ctx := context.Background()
			ev, _ := r.streamer.GetEvent(ctx)
			fmt.Println(ev.Event)
		}
	}()
	return
}
func (r *Input) doHandle() {
	for {
		logrus.Debug("get event")
		ctx := context.Background()
		e, er := r.streamer.GetEvent(ctx)
		if er != nil {
			logrus.Error(er)
			continue
		}
		//logrus.Debug(string(e.RawData))
		logrus.Debug(e.Header)
		//logrus.Debug(e.Event)
		msg, err := inputMessage(e)
		if err != nil {
			panic(err)
		}
		r.Ch <- msg
	}
	return
}

func (r *Input) handle() (err error) {
	go r.doHandle()
	return
}

func (r *Input) DataLine() chan *message.Message {
	return r.Ch
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
	r.Ch = make(chan *message.Message, 100000)
	return
}
