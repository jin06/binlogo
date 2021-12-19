package pipeline

// RecordPosition mysql replication position with pre position
type RecordPosition struct {
	PipelineName   string    `json:"pipeline_name"`
	Pre            *Position `json:"pre"`
	Now            *Position `json:"now"`
}

// OptionRecord record position options
type OptionRecord func(record *RecordPosition)

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
