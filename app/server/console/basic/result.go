package basic

import "fmt"

type Result struct {
	Code Code        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) Result {
	return Result{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	}
}

func Fail(msg interface{}) Result {
	errMsg := fmt.Sprintf("%v", msg)
	result := Result{
		Code: CodeUnknownError,
		Msg:  errMsg,
	}
	return result
}

func FailCode(code Code) Result {
	return Result{
		Code: code,
		Msg:  code.String(),
	}
}
