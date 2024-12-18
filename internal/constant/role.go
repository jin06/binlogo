package constant

type Role string

const (
	// LEADER node is leader
	LEADER Role = "leader"
	// FOLLOWER node is follower
	FOLLOWER Role = "follower"
)
