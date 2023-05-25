package routers

import (
    "github.com/gin-gonic/gin"
    "poop.fi/poop-server/internal/server/api/whitelist"
)

func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    {
        apiWhitelist := r.Group("/poop/whitelist")
        apiWhitelist.POST("/info", whitelist.Info)
    }

    return r
}
