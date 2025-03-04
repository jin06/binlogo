package constant

import (
	"fmt"
)

const (
	CodePipelineInstanceExists = 1001 + iota
)

var (
	// the error occurred during the registration of the instance key
	ErrPipelineInstanceExists = NewErr(CodePipelineInstanceExists, "failed to register pipeline instance, instance already exists")
)

func NewErr(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code:%d Message:%s", e.Code, e.Message)
}
