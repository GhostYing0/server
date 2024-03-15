package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func (g *Gin) ResponseSuc(msgs ...string) {
	var msgRet string
	data := make(map[string]interface{})
	for _, msg := range msgs {
		msgRet += " " + msg
	}

	g.C.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  msgRet,
		Data: data,
	})
	return
}

func (g *Gin) ResponseSucMsg(data interface{}, msgs ...string) {
	var msgRet string

	for _, msg := range msgs {
		msgRet += " " + msg
	}

	g.C.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  msgRet,
		Data: data,
	})
	return
}

func (g *Gin) ResponseErr(msgs ...string) {
	msgRet := "服务器发生错误"
	data := make(map[string]interface{})
	for _, msg := range msgs {
		msg += " "
		msgRet += msg
	}

	g.C.JSON(http.StatusOK, Response{
		Code: http.StatusInternalServerError,
		Msg:  msgRet,
		Data: data,
	})
	return
}
