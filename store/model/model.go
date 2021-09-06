package model

type Model interface {
	Val() string
	Key() string
	Unmarshal([]byte) error
}

