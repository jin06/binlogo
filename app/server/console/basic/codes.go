package basic

import "strconv"

type Code int

func (c Code) String() string {
	if val, ok := CODES[c]; ok {
		return val
	}
	return "unknown code " + strconv.Itoa(int(c))
}

const (
	CodeSuccess               Code = 20000
	CodeParamsFormat          Code = 40001
	CodeParamsRequired        Code = 40002
	CodeUnknownError          Code = 500
	CodeRequestDBFail         Code = 50001
	CodeLoginFail             Code = 50008
	CodeTokenExpired          Code = 50014
	CodeBadUsernameOrPassword Code = 60204
)

var CODES = map[Code]string{
	CodeSuccess:               "Success",
	CodeParamsFormat:          "Params error, json format",
	CodeParamsRequired:        "Required params is empty",
	CodeRequestDBFail:         "Search DB error",
	CodeLoginFail:             "Login failed, unable to get user details.",
	CodeTokenExpired:          "Token expired",
	CodeBadUsernameOrPassword: "Account and password are incorrect.",
	CodeUnknownError:          "Unknown error",
}
