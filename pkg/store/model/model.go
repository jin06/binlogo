package model

// Model base model
type Model interface {
	Val() string
	Key() string
	Unmarshal([]byte) error
}

type Values map[string]any
