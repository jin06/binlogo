package handler

import "fmt"

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) Result {
	return Result{
		Code: 20000,
		Msg:  "success",
		Data: data,
	}
}

func Fail(msg interface{}) Result {
	errMsg := fmt.Sprintf("%v", msg)
	result := Result{
		Code: 500,
		Msg: errMsg,
	}
	return result
}
