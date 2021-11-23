package model

// Header will deprecated
type Header struct {
	Revision int64
}

// GetHeader will deprecated
func (h *Header) GetHeader() *Header {
	return h
}

// SetHeader will deprecated
func (h *Header) SetHeader(n *Header) {
	*h = *n
}
