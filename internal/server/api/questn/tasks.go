package questn

import (
    "fmt"
    "github.com/gin-gonic/gin"
    ShitcrycleHistory "poop.fi/poop-server/internal/service/shitcrycle_history"
    UserChainInfo "poop.fi/poop-server/internal/service/user_chain_info"
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

func Task2(c *gin.Context) {
    addressStr := c.Query("address")
    userChainInfo, err := UserChainInfo.GetByChainAndAddress("bscmain", addressStr)
    if err != nil || userChainInfo == nil {
        c.String(200, fmt.Sprintf("{\"error\":{\"code\":0,\"message\":\"ok\"},\"data\":{\"result\":false}}"))
    } else if userChainInfo.Referral == nil {
        c.String(200, fmt.Sprintf("{\"error\":{\"code\":0,\"message\":\"ok\"},\"data\":{\"result\":false}}"))
    } else {
        c.String(200, fmt.Sprintf("{\"error\":{\"code\":0,\"message\":\"ok\"},\"data\":{\"result\":true}}"))
    }
}
