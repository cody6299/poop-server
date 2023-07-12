package routers

import (
    "github.com/gin-gonic/gin"
    "poop.fi/poop-server/internal/server/api/whitelist"
    "poop.fi/poop-server/internal/server/api/referral"
    "poop.fi/poop-server/internal/server/api/price"
    "poop.fi/poop-server/internal/server/api/questn"
)

func InitRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    {
        apiWhitelist := r.Group("/poop/whitelist")
        apiWhitelist.POST("/info", whitelist.Info)
    }

    {
        apiReferral := r.Group("/poop/referral")
        apiReferral.POST("/address", referral.Address)
        apiReferral.POST("/code", referral.Code)
        apiReferral.POST("/reward", referral.Reward)
        apiReferral.GET("/reward-history", referral.RewardHistory)
    }

    {
        apiPrice := r.Group("/poop/price")
        apiPrice.POST("/history", price.History)
        apiPrice.POST("/all", price.All)
    }

    {
        apiQuestN := r.Group("/poop/questn")
        apiQuestN.GET("/task1", questn.Task1)
        apiQuestN.GET("/task2", questn.Task2)
    }

    return r
}
