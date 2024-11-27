package pipeline

import (
	"encoding/json"
)

// Filter of pipeline
type Filter struct {
	Type FilterType `json:"type" redis:"type"`
	Rule string     `json:"rule" redis:"rule"`
}

func (m *Filter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *Filter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type Filters []Filter

func (m Filters) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m Filters) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// FilterType types of filter
type FilterType string

const (
	// FILTER_WHITE white list
	FILTER_WHITE FilterType = "white"
	// FILTER_BLACK black list
	FILTER_BLACK FilterType = "black"
)

// BlackFilter returns a black filter
func BlackFilter(rule string) (f *Filter) {
	f = &Filter{
		Type: FILTER_BLACK,
		Rule: rule,
	}
	return
}

// WhiteFilter returns a white filter
func WhiteFilter(rule string) (f *Filter) {
	f = &Filter{
		Type: FILTER_WHITE,
		Rule: rule,
	}
	return
}
