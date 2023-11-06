package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

type ErrorWithCode interface {
	ErrorCode() int
}
type ErrorWithHttpCode interface {
	ErrorHttpCode() int
}

// ResultErr 处理错误，如果错误为nil，则返回成功，否则按照错误类型返回
func ResultErr(data interface{}, e error, c *gin.Context) {
	var httpCode = http.StatusOK
	var code = ERROR
	var msg = "内部错误"

	if e == nil {
		code, msg = SUCCESS, "操作成功"
	} else {
		msg = e.Error()
		if ex, ok := e.(ErrorWithCode); ok {
			code = ex.ErrorCode()
		}
		if ex, ok := e.(ErrorWithHttpCode); ok && ex.ErrorHttpCode() != 0 {
			httpCode = ex.ErrorHttpCode()
		}
	}
	c.JSON(httpCode, Response{code, data, msg})
}

func Err(e error, c *gin.Context) {
	ResultErr(nil, e, c)
}
