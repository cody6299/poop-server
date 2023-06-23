package price

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/server/model"
    "poop.fi/poop-server/internal/config"
    "poop.fi/poop-server/internal/utils"
    PriceInfo "poop.fi/poop-server/internal/service/price_info"
)

type HistoryParam struct {
    ChainId uint    `json:"chain_id" binding:"required"`
    Type    string  `json:"type" binding:"required"`
    Size    uint    `json:"size" binding:"required"`
}

type PriceRecordResponse struct {
    BeginTime   uint64  `json:"begin_time"`
    EndTime     uint64  `json:"end_timeee"`
    PriceOpen   uint64  `json:"price_open"`
    PriceHigh   uint64  `json:"price_high"`
    PriceLow    uint64  `json:"price_low"`
    PriceClose  uint64  `json:"price_close"`
}

type HistoryResponse struct {
    Records    []*PriceRecordResponse `json:"price_record"`
}

func getPriceType(t string) *string {
    switch(t) {
    case "1minute":
        return utils.PTR(fmt.Sprintf("price_%d", 60))
    case "5minute":
        return utils.PTR(fmt.Sprintf("price_%d", 5 * 60))
    case "15minute":
        return utils.PTR(fmt.Sprintf("price_%d", 15 * 60))
    case "30minute":
        return utils.PTR(fmt.Sprintf("price_%d", 30 * 60))
    case "1hour":
        return utils.PTR(fmt.Sprintf("price_%d", 60 * 60))
    case "4hour":
        return utils.PTR(fmt.Sprintf("price_%d", 4 * 60 * 60))
    case "1day":
        return utils.PTR(fmt.Sprintf("price_%d", 24 * 60 * 60))
    case "3day":
        return utils.PTR(fmt.Sprintf("price_%d", 3 * 24 * 60 * 60))
    case "7day":
        return utils.PTR(fmt.Sprintf("price_%d", 7 * 24 * 60 * 60))
    case "14day":
        return utils.PTR(fmt.Sprintf("price_%d", 14 * 24 * 60 * 60))
    case "1month":
        return utils.PTR(fmt.Sprintf("price_%d", 30 * 24 * 60 * 60))
    case "3month":
        return utils.PTR(fmt.Sprintf("price_%d", 3 * 30 * 24 * 60 * 60))
    case "6month":
        return utils.PTR(fmt.Sprintf("price_%d", 6 * 30 * 24 * 60 * 60))
    default:
        return nil
    }
}

func History(c *gin.Context) {
    g := model.Gin{C: c}

    param := &HistoryParam{}
    c.BindJSON(param)

    if param.Size > 10000 {
        g.Response(http.StatusOK, -2, "illegal price type", nil)
        return
    }

    chainName := config.GetChainNameByChainId(param.ChainId)
    if chainName == nil {
        g.Response(http.StatusOK, -2, "illegal chainId", nil)
        return
    }

    priceType := getPriceType(param.Type)
    if priceType == nil {
        g.Response(http.StatusOK, -2, "illegal price type", nil)
        return
    }

    priceInfoList, err := PriceInfo.GetRangeByChainAndType(*chainName, *priceType, 0, param.Size)
    if err != nil {
        log.Errorf("PriceInfo.GetRangeByChainAndType failed: %v", err)
        return
    }

    records := []*PriceRecordResponse{}
    for _, priceInfo := range *priceInfoList {
        record := &PriceRecordResponse{
            BeginTime: priceInfo.BeginTime,
            EndTime: priceInfo.EndTime,
            PriceOpen: priceInfo.PriceOpen,
            PriceHigh: priceInfo.PriceHigh,
            PriceLow: priceInfo.PriceLow,
            PriceClose: priceInfo.PriceClose,
        }
        records = append(records, record)
    }

    response := &HistoryResponse{
        Records: records,
    }

    g.Response(http.StatusOK, 0, "ok", response)
}

type AllParam struct {
    ChainId uint    `json:"chain_id" binding:"required"`
}

type AllResponse struct {
    Records    []*PriceRecordResponse `json:"price_record"`
    Interval   uint                   `json:"interval"`
}

func All(c *gin.Context) {
    g := model.Gin{C: c}

    param := &AllParam{}
    c.BindJSON(param)

    chainName := config.GetChainNameByChainId(param.ChainId)
    if chainName == nil {
        g.Response(http.StatusOK, -2, "illegal chainId", nil)
        return
    }

    priceType := getPriceType("5minute")
    if priceType == nil {
        g.Response(http.StatusOK, -2, "illegal price type", nil)
        return
    }

    priceInfoList, err := PriceInfo.GetByChainAndType(*chainName, *priceType)
    if err != nil {
        log.Errorf("PriceInfo.GetRangeByChainAndType failed: %v", err)
        return
    }

    records := []*PriceRecordResponse{}
    for _, priceInfo := range *priceInfoList {
        record := &PriceRecordResponse{
            BeginTime: priceInfo.BeginTime,
            EndTime: priceInfo.EndTime,
            PriceOpen: priceInfo.PriceOpen,
            PriceHigh: priceInfo.PriceHigh,
            PriceLow: priceInfo.PriceLow,
            PriceClose: priceInfo.PriceClose,
        }
        records = append(records, record)
    }

    response := &AllResponse{
        Records: records,
        Interval: 60,
    }

    g.Response(http.StatusOK, 0, "ok", response)
}
