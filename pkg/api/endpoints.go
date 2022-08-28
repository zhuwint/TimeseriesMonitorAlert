package api

type ErrorCode string

const (
	prefix           ErrorCode = "error."
	RequestBodyError ErrorCode = prefix + "40000"
	QueryParamError  ErrorCode = prefix + "40010"
)

type ReplyError struct {
	Code ErrorCode `json:"code"`
	Msg  string    `json:"msg"`
}
