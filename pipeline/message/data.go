package message

const TYPE_INSERT = "insert"

type Insert struct {
	Old map[string]string
	New map[string]string
}

type Update struct {
	Old map[string]interface{} `json:"old"`
	New map[string]interface{} `json:"new"`

}
