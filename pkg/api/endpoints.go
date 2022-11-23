package api

type ErrorCode string

const (
	prefix           ErrorCode = "error."
	RequestBodyError ErrorCode = prefix + "40000"
	QueryParamError  ErrorCode = prefix + "40010"
	InternelError    ErrorCode = prefix + "40020"
	ResourceNotFound ErrorCode = prefix + "40030"
)

type ReplyError struct {
	Code ErrorCode `json:"code"`
	Msg  string    `json:"msg"`
}

type ReplyJson struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
