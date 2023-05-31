package errorx

import (
	"fmt"
)

type CodeError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

// 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.Code
}

// 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.Msg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.Code, e.Msg)
}

func New(code uint32, errMsg string) *CodeError {
	return &CodeError{Code: code, Msg: errMsg}
}

func NewMsg(errMsg string) *CodeError {
	return &CodeError{Code: ERROR, Msg: errMsg}
}

func NewMsgF(msgFormat string, args ...interface{}) *CodeError {
	return &CodeError{Code: ERROR, Msg: fmt.Sprintf(msgFormat, args)}
}

func NewCodeMsgF(code uint32, msgFormat string, args ...interface{}) *CodeError {
	return &CodeError{Code: code, Msg: fmt.Sprintf(msgFormat, args)}
}

func (e *CodeError) ToJson() string {
	return fmt.Sprintf("{\"code\": %d, \"msg\": \"%s\"}", e.Code, e.Msg)
}
