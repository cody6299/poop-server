package whitelist

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "poop.fi/poop-server/internal/server/model"
)

type InfoParam struct {
    Address string `json:"address" binding:"required"`
}

func Info(c *gin.Context) {
    g := model.Gin{C: c}

    param := &InfoParam{}
    c.BindJSON(param)
    fmt.Printf("param: %v\n", param)

    g.Response(http.StatusOK, 0, "ok", nil)
}
