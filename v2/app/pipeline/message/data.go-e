package message

// MessageType for pipeline handle
type MessageType byte

var (
	TYPE_EMPTY        MessageType = 0
	TYPE_INSERT       MessageType = 1
	TYPE_UPDATE       MessageType = 2
	TYPE_DELETE       MessageType = 3
	TYPE_CREATE_TABLE MessageType = 4
	TYPE_ALTER_TABLE  MessageType = 5
)

// String returns MessageType's string
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
	case TYPE_EMPTY:
		{
			return "empty"
		}
	}
	return ""
}

// Insert for mysql insert
type Insert struct {
	New map[string]interface{} `json:"new"`
}

// Update for mysql update
type Update struct {
	Old map[string]interface{} `json:"old"`
	New map[string]interface{} `json:"new"`
}

// Delete for mysql delete
type Delete struct {
	Old map[string]interface{} `json:"old"`
}

// CreateTable for mysql ddl
type CreateTable struct {
}

// AlterTable for mysql ddl
type AlterTable struct {
}
