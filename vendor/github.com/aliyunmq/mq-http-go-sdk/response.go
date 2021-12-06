package mq_http_sdk

type MessageResponder interface {
	Responses() []*MessageResponse
}

func (r *PublishMessageResponse) Responses() []*MessageResponse {
	var responseSlice []*MessageResponse
	responseSlice = append(responseSlice, &r.MessageResponse)
	return responseSlice
}

func (r *ConsumeMessageResponse) Responses() []*MessageResponse {
	var responseSlice []*MessageResponse
	for i := range r.Messages {
		responseSlice = append(responseSlice, &r.Messages[i].MessageResponse)
	}
	return responseSlice
}