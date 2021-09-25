package message

type MessageType byte

var (
	TYPE_INSERT       MessageType = 1
	TYPE_UPDATE       MessageType = 2
	TYPE_DELETE       MessageType = 3
	TYPE_CREATE_TABLE MessageType = 4
	TYPE_ALTER_TABLE  MessageType = 5
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
	case TYPE_CREATE_TABLE:
		{
			return "create_table"
		}
	case TYPE_ALTER_TABLE:
		{
			return "alter_table"
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

type CreateTable struct {
}

type AlterTable struct {
}
