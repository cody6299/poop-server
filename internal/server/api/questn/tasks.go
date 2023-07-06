package questn

import (
    "fmt"
    "github.com/gin-gonic/gin"
    ShitcrycleHistory "poop.fi/poop-server/internal/service/shitcrycle_history"
)

func Task1(c *gin.Context) {
    addressStr := c.Query("address")

    num, err := ShitcrycleHistory.CountByAddress(addressStr)
    if err != nil || num == 0{
        c.String(200, fmt.Sprintf("{\"error\":{\"code\":0,\"message\":\"ok\"},\"data\":{\"result\":false}}"))
    } else {
        c.String(200, fmt.Sprintf("{\"error\":{\"code\":0,\"message\":\"ok\"},\"data\":{\"result\":true}}"))
    }
}
