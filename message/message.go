package message

import "time"

type BinlogMessage interface {
	Serialize() string
}

type Message struct {
	Header struct{
		Date time.Time
		LogPosition uint32
	}
	Body struct{}
}

func Serialize(m *Message) (ret string){
	return
}





