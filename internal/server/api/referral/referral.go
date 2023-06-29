package referral

import (
    "fmt"
    "math/big"
    "bytes"
    "regexp"
    "strconv"
    "net/http"
    "encoding/csv"
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/server/model"
    UserReferralCode "poop.fi/poop-server/internal/service/user_referral_code"
    UserChainInfo "poop.fi/poop-server/internal/service/user_chain_info"
    ReferralRewardInfo "poop.fi/poop-server/internal/service/referral_reward_info"
    "poop.fi/poop-server/internal/utils"
    "github.com/ethereum/go-ethereum/common"
    "poop.fi/poop-server/internal/config"
)

type AddressParam struct {
    Address string `json:"address" binding:"required"`
}

type AddressResponse struct {
    Address         string `json:"address"`
    ReferralCode    string `json:"referral_code"`
}

func Address(c *gin.Context) {
    g := model.Gin{C: c}

    param := &AddressParam{}
    c.BindJSON(param)

    re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
    addressValid := re.MatchString(param.Address)
    if !addressValid {
        g.Response(http.StatusOK, -2, "illegal address", nil)
        return
    }
    address := common.HexToAddress(param.Address)

    //首先从数据库中查询是否存在
    userReferralCode, err := UserReferralCode.GetByAddress(address.Hex()) 
    if err != nil {
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    if userReferralCode == nil {
        userReferralCode = &UserReferralCode.UserReferralCode{
            Address: address.Hex(),
        }
        UserReferralCode.Save(userReferralCode)
        //无论成功与否都继续
        userReferralCode, err := UserReferralCode.GetByAddress(address.Hex())
        if err != nil || userReferralCode == nil {
            //如果还是查询不到,那就是出错了 
            log.Errorf("get user referral code for %s failed err=%v result=%v", address.Hex(), err, userReferralCode)
            g.Response(http.StatusOK, -1, "system error", nil)
            return
        }
    }
    //这里我们找到了最新的UserReferralCode,根据他的ID就可以算出来对应的CODE
    if userReferralCode == nil || userReferralCode.Id < 123456 {
        log.Errorf("get user referral code for %s failed err=%v result=%v", address.Hex(), err, userReferralCode)
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    if userReferralCode.ReferralCode == nil {
        userReferralCode.ReferralCode = utils.PTR(UserReferralCode.GenerateReferralCode(userReferralCode.Id))
        //再更新ReferralCode到DB
        _, err = UserReferralCode.UpdateReferralCode(userReferralCode)
        if err != nil {
            log.Errorf("update referral code for %s failed %v", address.Hex(), err)
            g.Response(http.StatusOK, -1, "system error", nil)
            return
        }
    }

    response := AddressResponse {
        Address: userReferralCode.Address,
        ReferralCode: *userReferralCode.ReferralCode,
    }

    g.Response(http.StatusOK, 0, "ok", response)
}

type CodeParam struct {
    Code string `json:"code" binding:"required"`
}

type CodeResponse struct {
    Address         string `json:"address"`
    ReferralCode    string `json:"referral_code"`
}

func Code(c *gin.Context) {
    g := model.Gin{C: c}

    param := &CodeParam{}
    c.BindJSON(param)

    //首先从数据库中查询是否存在
    userReferralCode, err := UserReferralCode.GetByCode(param.Code) 
    if err != nil {
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    if userReferralCode == nil {
        g.Response(http.StatusOK, -3, "not exist", nil)
        return
    }

    if userReferralCode.ReferralCode == nil {
        userReferralCode.ReferralCode = utils.PTR(UserReferralCode.GenerateReferralCode(userReferralCode.Id))
        //再更新ReferralCode到DB
        _, err = UserReferralCode.UpdateReferralCode(userReferralCode)
        if err != nil {
            log.Errorf("update referral code for %s failed %v", userReferralCode.Address, err)
            g.Response(http.StatusOK, -1, "system error", nil)
            return
        }
    }
    
    response := &CodeResponse {
        Address: userReferralCode.Address,
        ReferralCode: *userReferralCode.ReferralCode,
    }

    g.Response(http.StatusOK, 0, "ok", response)
}

type RewardParam struct {
    ChainId     uint    `json:"chain_id" binding:"required"`
    Address     string  `json:"address" binding:"required"`
    Page        uint    `json:"page" binding:"required"`
    Size        uint    `json:"size" binding:"required"`
}

type UserInfoResponse struct {
    UserId          *uint64     `json:"user_id"`
    Referral        *string     `json:"referral"`
    ReferralTime    *int64      `json:"referral_time"`
    InviteNum       uint        `json:"invite_num"`
    InviteReward    uint64      `json:"invite_reward"`
    RewardNum       uint        `json:"reward_num"`
}

type PageInfoResponse struct {
    TotalNum    uint    `json:"total_num"`
    TotalPage   uint    `json:"total_page"`
    Page        uint    `json:"page"`
    Size        uint    `json:"size"`
}

type RewardRecordResponse struct {
    Address     string  `json:"address"`
    Amount      uint64  `json:"amount"`
    Time        int64   `json:"time"`
    TxHash      string  `json:"tx_hash"`
}

type RewardResponse struct {
    Chain       string                      `json:"chain"`
    Address     string                      `json:"address"`
    UserInfo    *UserInfoResponse           `json:"user_info"`
    PageInfo    *PageInfoResponse           `json:"page_info"`
    Records     []*RewardRecordResponse     `json:"records"`
}

func Reward(c *gin.Context) {
    g := model.Gin{C: c}

    param := &RewardParam{}
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

    if param.Page < 1 {
        g.Response(http.StatusOK, -2, "illegal page", nil)
        return
    }
    if param.Size > 100 {
        g.Response(http.StatusOK, -2, "illegal size", nil)
        return
    }

    //首先获取UserChainInfo
    userChainInfo, err := UserChainInfo.GetByChainAndAddress(*chainName, address)
    if err != nil {
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    if userChainInfo == nil {
        userChainInfo = &UserChainInfo.UserChainInfo {
            Chain: *chainName,
            Address: address,
            InviteNum: 0,
            InviteReward: 0,
            RewardNum: 0,
        }
    }

    userInfoResponse := &UserInfoResponse{
        UserId : userChainInfo.UserId,
        Referral: userChainInfo.Referral,
        InviteNum: userChainInfo.InviteNum,
        InviteReward: userChainInfo.InviteReward,
        RewardNum: userChainInfo.RewardNum,
    }
    if userChainInfo.ReferralTime == nil {
        userInfoResponse.ReferralTime = nil
    } else {
        userInfoResponse.ReferralTime = utils.PTR(userChainInfo.ReferralTime.Unix())
    }
    
    pageInfo := &PageInfoResponse{
        TotalNum: userChainInfo.RewardNum,
        Page: param.Page,
        Size: param.Size,
    }
    if pageInfo.TotalNum % pageInfo.Size != 0 {
        pageInfo.TotalPage = pageInfo.TotalNum / pageInfo.Size + 1
    } else {
        pageInfo.TotalPage = pageInfo.TotalNum / pageInfo.Size
    }

    records := []*RewardRecordResponse{}
    offset := (pageInfo.Page - 1) * pageInfo.Size
    referralRewardInfoList, err := ReferralRewardInfo.GetRangeByChainAndAddress(*chainName, address, offset, pageInfo.Size)
    if err != nil {
        log.Errorf("ReferralRewardInfo.GetRangeByChainAndAddress error: %v", err)
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    for _, info := range *referralRewardInfoList {
        record := &RewardRecordResponse {
            Address: info.InviteAddress,
            Amount: info.RewardAmount,
            Time: info.RewardTime.Unix(),
            TxHash: info.TxHash,
        }
        records = append(records, record)
    }

    response := &RewardResponse{
        Chain: *chainName,
        Address: address,
        UserInfo: userInfoResponse,
        PageInfo: pageInfo,
        Records: records,
    }

    g.Response(http.StatusOK, 0, "ok", response)
}

type RewardHistoryParam struct {
    ChainId     uint    `json:"chain_id" binding:"required"`
    StartTime   uint    `json:"start_time"`
    EndTime     uint    `json:"end_time"`
    SortByTag   string  `json:"sort_by_tag"`
    Size        uint    `json:"size"`
}

func RewardHistory(c *gin.Context) {
    g := model.Gin{C: c}

    chainIdStr := c.Query("chain_id")
    chainId, err := strconv.Atoi(chainIdStr)
    if err != nil {
        g.Response(http.StatusOK, -1, "illegal chainId", nil)
        return
    }
    chainName := config.GetChainNameByChainId(uint(chainId))
    if chainName == nil {
        fmt.Println("test124")
        g.Response(http.StatusOK, -1, "illegal chainId", nil)
        return
    }

    startTimeStr := c.DefaultQuery("start_time", "0")
    startTime, err := strconv.Atoi(startTimeStr)
    if err != nil {
        g.Response(http.StatusOK, -1, "illegal start_time", nil)
        return
    }

    endTimeStr := c.DefaultQuery("end_time", "")
    endTime, err := strconv.Atoi(endTimeStr)
    if err != nil {
        g.Response(http.StatusOK, -1, "illegal end_time", nil)
        return
    }

    start := 0
    records := []ReferralRewardInfo.AggregationRecord{}
    referralRewardInfoList, err := ReferralRewardInfo.AggregationByTime(*chainName, startTime, endTime)
    if err != nil {
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    if (len(*referralRewardInfoList) > 0) {
        records = append(records, (*referralRewardInfoList)...)
        start = start + len(*referralRewardInfoList)
    }
    log.Debugf("records: %v", len(records))

    start = 0
    userChainInfoRecords := []UserChainInfo.AggregationRecord{}
    userChainInfoList, err := UserChainInfo.AggregationByTime(*chainName, startTime, endTime)
    if err != nil {
        g.Response(http.StatusOK, -1, "system error", nil)
        return
    }
    if len(*userChainInfoList) > 0 {
        userChainInfoRecords = append(userChainInfoRecords, (*userChainInfoList)...)
        start = start + len(*userChainInfoList)
    }
    log.Debugf("records: %v", len(userChainInfoRecords))
    var userChainInfoMap map[string]*UserChainInfo.AggregationRecord = make(map[string]*UserChainInfo.AggregationRecord)
    for _, record1 := range userChainInfoRecords {
        userChainInfoMap[*(record1.Referral)] = &record1
    }
    //fmt.Printf("address,total_reward_amount,total_reward_num,total_referral_num")

    dataBytes := new(bytes.Buffer)
    dataBytes.WriteString("\xEF\xBB\xBF")
    writer := csv.NewWriter(dataBytes)
    writer.Write([]string{"id", "address","total_reward_amount","total_reward_num","total_referral_num"})
    units, _ := new (big.Float).SetString("1000000000000000000") 
    for index, record2 := range records {
        address := record2.Address 
        referralNum := uint64(0)
        if userChainInfoMap[address] != nil {
            referralNum = userChainInfoMap[address].TotalReferralNum 
        }
        totalRewardAmount := new(big.Float).SetInt(record2.TotalRewardAmount)
        bodyList := []string{
            strconv.FormatUint(uint64(index + 1), 10),
            address, 
            fmt.Sprintf("%v BNB", totalRewardAmount.Quo(totalRewardAmount, units)),
            strconv.FormatUint(record2.TotalRewardNum, 10),
            strconv.FormatUint(referralNum, 10),
        }
        writer.Write(bodyList)
    }
    writer.Flush()

    c.Writer.Header().Set("Content-type", "application/octet-stream")
    c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", "referral.csv"))
    c.String(200, dataBytes.String())
    //g.Response(http.StatusOK, 0, "ok", nil)
}
