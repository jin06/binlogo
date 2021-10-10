package role

type Role byte

const (
	LEADER   Role = 1
	FOLLOWER Role = 2
)

func (r Role) String() (str string) {
	switch r {
	case LEADER:
		{
			str = "leader"
		}
	case FOLLOWER:
		{
			str = "follower"
		}
	default:
		{
			str = "follower"
		}
	}
	return
}
