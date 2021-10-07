package routers

type Result struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func NewSuccess()  *Result {
	ret := &Result{
		Code : 200,
		Message :"ok",
	}
	return ret
}
