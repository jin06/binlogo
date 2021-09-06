package model

type Model interface {
	Val() string
	Key() string
	Unmarshal([]byte) error
}

//type ETCD struct {
//}
//
//func (etcd *ETCD) Val() (val string) {
//	b, _ := json.Marshal(etcd)
//	val = string(b)
//	return
//}
