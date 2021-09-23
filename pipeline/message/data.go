package message

type MessageType byte

var (
	TYPE_INSERT MessageType = 1
	TYPE_UPDATE MessageType = 2
	TYPE_DELETE MessageType = 3
)

func (mt MessageType) String() string {
	switch mt {
	case TYPE_INSERT:
		{
			return "insert"
		}
	case TYPE_UPDATE:
		{
			return "update"
		}
	case TYPE_DELETE:
		{
			return "delete"
		}
	}
	return ""
}

type Insert struct {
	New map[string]interface{} `json:"new"`
}

type Update struct {
	Old map[string]interface{} `json:"old"`
	New map[string]interface{} `json:"new"`
}

type Delete struct {
	Old map[string]interface{} `json:"old"`
}
