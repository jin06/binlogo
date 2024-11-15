package node

import (
	"encoding/json"
)

// Role store struct
type Role struct {
	Master bool `json:"master" redis:"master"`
	Admin  bool `json:"admin" redis:"admin"`
	Worker bool `json:"worker" redis:"worker"`
}

func (r Role) MarshalBinary() (data []byte, err error) {
	// encoding.BinaryMarshaler
	return json.Marshal(r)
}

func (r Role) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}
