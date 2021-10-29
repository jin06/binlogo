package pipeline

type Filter struct {
	Type FilterType `json:"type"`
	Rule string     `json:"rule"`
}

type FilterType string

const (
	FILTER_WHITE FilterType = "white"
	FILTER_BLACK FilterType = "black"
)
