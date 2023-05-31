package result

import (
	"fmt"
	"koala/gateway/internal/tools/errorx"
	"net/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

// HttpResult 返回
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		HttpResultOk(w, resp)
	} else {
		HttpResultError(w, err)
	}
}

// HttpResultOk 成功返回
func HttpResultOk(w http.ResponseWriter, resp interface{}) {
	res := Success(resp)
	httpx.WriteJson(w, http.StatusOK, res)
}

// HttpResultError 错误返回
func HttpResultError(w http.ResponseWriter, err error) {
	//错误返回
	errCode := errorx.ERROR
	errMsg := "服务器开小差啦，稍后再来试一试"

	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*errorx.CodeError); ok {
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
	} else {
		if gstatus, ok := status.FromError(causeErr); ok {
			grpcCode := uint32(gstatus.Code())
			errCode = grpcCode
			errMsg = gstatus.Message()
		}
	}

	httpx.WriteJson(w, http.StatusOK, Error(errCode, errMsg))
}

// ParamErrorResult 参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", errorx.RequestParamError, err.Error())
	httpx.WriteJson(w, http.StatusOK, Error(errorx.Param, errMsg))
}
