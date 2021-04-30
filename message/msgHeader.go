package message

import "time"

type MsgHeader struct {
	Date time.Time
	LogPosition uint32
	Type string
	Database string
	Table string
}
