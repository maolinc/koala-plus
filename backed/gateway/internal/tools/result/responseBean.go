package result

type ResponseSuccessBean struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"result"`
}
type NullJson struct{}

func Success(data interface{}) *ResponseSuccessBean {
	return &ResponseSuccessBean{0, "OK", data}
}

type ResponseErrorBean struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func Error(errCode uint32, errMsg string) *ResponseErrorBean {
	return &ResponseErrorBean{errCode, errMsg}
}
