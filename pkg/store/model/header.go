package model

type Header struct {
	Revision int64
}

func (h *Header) GetHeader() *Header {
	return h
}
func (h *Header) SetHeader(n *Header) {
	*h = *n
}
