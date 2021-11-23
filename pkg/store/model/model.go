package model

// Model base model
type Model interface {
	Val() string
	Key() string
	Unmarshal([]byte) error
}

// ModelH will deprecated
type ModelH interface {
	Model
	GetHeader() *Header
	SetHeader(*Header)
}
