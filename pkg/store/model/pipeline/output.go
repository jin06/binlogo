package pipeline

// Output pipeline output
type Output struct {
	Sender Sender `json:"sender" redis:"sender"`
}
