package input

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	message2 "github.com/jin06/binlogo/app/pipeline/message"
	"github.com/jin06/binlogo/pkg/store/model/pipeline"
	"github.com/sirupsen/logrus"
)

type Input struct {
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer
	OutChan  chan *message2.Message
	Options  *Options
}

func (r *Input) Run(ctx context.Context) (err error) {
	if err = r.connect(ctx); err != nil {
		return
	}
	go func() {
		for  {
			select {
				case <- ctx.Done():{
					r.syncer.Close()
					return
				}
			default:
				r.doHandle(ctx)
			}
		}
	}()
	return
}

func (r *Input) connect(ctx context.Context) (err error) {
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
	//r.syncer.GetNextPosition()
	if err != nil {
		return
	}
	r.streamer = streamer
	return
}

func (r *Input) doHandle(ctx context.Context) {
	logrus.Debug("Sync binlog from mysql")
	//ctx := context.Background()
	e, er := r.streamer.GetEvent(ctx)
	if er != nil {
		logrus.Error(er)
		return
	}
	logrus.Debug("Binlog event header : ", e.Header)
	msg, err := inputMessage(e)
	pos := r.syncer.GetNextPosition()
	fmt.Println(pos)
	msg.BinlogPosition = &pipeline.Position{
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

func (r *Input) preparePosition(pos *mysql.Position) (err error){
	if pos != nil {
		if pos.Name != "" {
			return
		}
	}
	return
}

func (r *Input) newestPosition() {
}


