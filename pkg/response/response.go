package response

import (
	"homework-system/pkg/errcode" //自定义错误包

	"github.com/gin-gonic/gin" //Gin框架
)

type Response struct {
	Code    int         `json:"code"`    //业务状态码
	Message string      `json:"message"` //响应消息
	Data    interface{} `json:"data"`    //响应数据
}

func Success(c *gin.Context, data interface{}) { //处理成功响应
	c.JSON(200, Response{ //调用JSON方法，返回HTTP状态码200和JSON数据
		Code:    errcode.Success,
		Message: errcode.Message[errcode.Success],
		Data:    data,
	})
}
func Error(c *gin.Context, httpCode int, code int) { //处理错误响应的函数
	c.JSON(httpCode, Response{
		Code:    code,
		Message: errcode.Message[code],
		Data:    nil,
	})
}
func ErrorWithMsg(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
