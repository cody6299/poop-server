package model

import (
    "github.com/gin-gonic/gin"
)

type Gin struct {
    C *gin.Context
}

type Response struct {
    Code    int         `json:"code"`
    Msg     string      `json:"message"`
    Data    interface{} `json:"data"`
}

func (g *Gin) Response(httpCode, code int, msg string, data interface{}) {
    g.C.JSON(httpCode, Response{
        Code: code,
        Msg: msg,
        Data: data,
    })
    return
}
