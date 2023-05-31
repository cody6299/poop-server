package whitelist

import (
    "regexp"
    "strings"
    "net/http"
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/server/model"
    "poop.fi/poop-server/internal/config"
    "github.com/ethereum/go-ethereum/common"
    WhitelistInfo "poop.fi/poop-server/internal/service/whitelist_info"
)

type InfoParam struct {
    ChainId uint `json:"chain_id" binding:"required"`
    Address string `json:"address" binding:"required"`
}

type InfoResponse struct {
    ChainName   string `json:"chain_id"`
    Address     string `json:"address"`
    MaxAmount   string `json:"max_amount"`
    Proof       []string `json:"proof"`
    
}

func Info(c *gin.Context) {
    g := model.Gin{C: c}

    param := &InfoParam{}
    c.BindJSON(param)

    re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
    addressValid := re.MatchString(param.Address)
    if !addressValid {
        g.Response(http.StatusOK, -2, "illegal address", nil)
        return
    }
    address := common.HexToAddress(param.Address).Hex()

    chainName := config.GetChainNameByChainId(param.ChainId)
    if chainName == nil {
        g.Response(http.StatusOK, -2, "illegal chainId", nil)
        return
    }

    whitelistInfo, err := WhitelistInfo.GetByChainAndAddress(*chainName, address)
    if err != nil {
        log.Errorf("WhitelistInfo.GetByChainAndAddress failed: %v", err)
    }
    if whitelistInfo == nil {
        g.Response(http.StatusOK, 0, "ok", nil)
        return
    }

    response := &InfoResponse{
        ChainName: *chainName,
        Address: address,
        MaxAmount: whitelistInfo.MaxAmount,
        Proof: strings.Split(whitelistInfo.Proof, ","),
    }

    g.Response(http.StatusOK, 0, "ok", response)
}
