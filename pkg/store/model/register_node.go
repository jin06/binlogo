package model

import "encoding/json"

type RegisterNode struct {
	Name string `json:"name" redis:"name"`
}

func (s *RegisterNode) Key() string {
	return "register_node"
}

func (s *RegisterNode) Val() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *RegisterNode) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}
