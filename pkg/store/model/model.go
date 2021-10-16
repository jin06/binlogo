package model

type Model interface {
	Val() string
	Key() string
	Unmarshal([]byte) error
}

type ModelH interface {
	Model
	GetHeader() *Header
	SetHeader(*Header)
}
