package pipeline

// RecordPosition mysql replication position with pre position
type RecordPosition struct {
	PipelineName string    `json:"pipeline_name" redis:"pipeline_name"`
	Pre          *Position `json:"pre" redis:"pre"`
	Now          *Position `json:"now" redis:"now"`
}

// NewRecordPosition return empty RecordPosition
func NewRecordPosition(opts ...OptionRecord) *RecordPosition {
	r := &RecordPosition{
		Pre: &Position{},
		Now: &Position{},
	}
	for _, v := range opts {
		v(r)
	}
	return r
}

// OptionRecord record position options
type OptionRecord func(record *RecordPosition)

// WithPipelineName set RecordPosition pipeline name
func WithPipelineName(name string) OptionRecord {
	return func(record *RecordPosition) {
		record.PipelineName = name
	}
}

// WithPre record position pre
func WithPre(p *Position) OptionRecord {
	return func(record *RecordPosition) {
		record.Pre = p
	}
}

// WithNow record position now
func WithNow(n *Position) OptionRecord {
	return func(record *RecordPosition) {
		record.Now = n
	}
}
