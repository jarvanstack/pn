package response

import (
	"dc3/public/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 错误信息
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Ok 返回成功
//	@msg:	返回消息
//	@data:	返回成功的数据
func Ok(c *gin.Context, data ...interface{}) {
	resp := Response{
		Code: 200,
		Msg:  "OK",
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	c.JSON(http.StatusOK, resp)
}

// Err 返回失败
//	@param:	httpCode 错误码
//	@param:	msg 错误消息
//	@param:	data 额外信息（可选）

func Err(c *gin.Context, err error, data ...interface{}) {

	resp := Response{}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	switch err.(type) {
	case errs.IErr:
		resp.Code = err.(errs.IErr).Code()
		resp.Msg = err.(errs.IErr).Error()
	default:
		resp.Code = http.StatusInternalServerError
		resp.Msg = err.Error()
	}

	c.JSON(http.StatusOK, resp)
}
