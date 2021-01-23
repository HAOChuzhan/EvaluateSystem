package app

import (
	"github.com/gin-gonic/gin"
)

//封装一个 ‘*gin.Context’
type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response 设置gin.JSON
func (g *Gin) Response(httpCode int, errCode string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  errCode,
		Data: data,
	})

	return
}

//Respose 设置无data
