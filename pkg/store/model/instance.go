package model

type InstanceRole string

const (
	INSTANCE_LEADER  InstanceRole = "leader"
	INSTANCE_FLOOWER InstanceRole = "follower"
)

type Instance struct {
	Role *InstanceRole `json:"role"`
}
