package message

type Message struct {
	Name string
	Content *Content
}

type Content struct {
	Head *Head
	Event *Event
}

type Head struct {
	Time uint32
}

type Event struct {
	Type EventType
}
