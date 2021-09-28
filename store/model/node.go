package model

import "encoding/json"

type Node struct {
	Name string `json:"name"`
}

func (s *Node) Key() (key string) {
	key = "nodes/" + s.Name
	return
}

func (s *Node) Val() (val string) {
	b, _ := json.Marshal(s)
	val = string(b)
	return
}
