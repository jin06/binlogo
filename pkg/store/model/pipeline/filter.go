package pipeline

// Filter of pipeline
type Filter struct {
	Type FilterType `json:"type"`
	Rule string     `json:"rule"`
}

// FilterType types of filter
type FilterType string

const (
	// FILTER_WHITE white list
	FILTER_WHITE FilterType = "white"
	// FILTER_BLACK black list
	FILTER_BLACK FilterType = "black"
)
