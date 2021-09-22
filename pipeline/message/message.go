package message

type Message struct {
	Content *Content
}

type Content struct {
	Head *Head       `json:"head"`
	Data interface{} `json:"data"`
}

type Head struct {
	Type     string `json:"type"`
	Time     uint32 `json:"time"`
	Database string `json:"database"`
	Table    string `json:"table"`
}
