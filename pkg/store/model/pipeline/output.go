package pipeline

// Output pipeline output
type Output struct {
	Sender Sender `json:"sender" redis:"sender"`
}

// func (m Output) MarshalBinary() (data []byte, err error) {
// 	return json.Marshal(m)
// }

// func (m Output) UnmarshalBinary(data []byte) error {
// 	return json.Unmarshal(data, &m)
// }

// func (m Output) UnmarshalText(text []byte) error {
// 	return json.Unmarshal(text, &m)
// }

// EmptyOutput return a new empty Output object
func EmptyOutput() *Output {
	return &Output{
		Sender: Sender{
			Name: "",
			Type: "",
		},
	}
}
