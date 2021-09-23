package message

const TYPE_INSERT = "insert"
const TYPE_UPDATE = "update"
const TYPE_DELETE = "delete"

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
