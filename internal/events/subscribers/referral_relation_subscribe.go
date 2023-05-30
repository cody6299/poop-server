package subscribers

import (
    "os"
    "fmt"
    "time"
    "errors"
    "strings"
    "math/big"
    "gorm.io/gorm"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/config"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/accounts/abi"
    UserChainInfo "poop.fi/poop-server/internal/service/user_chain_info"
    "poop.fi/poop-server/internal/utils"
    
)

type ReferralRelationSubscribe struct {
    config *config.Config
    chain *config.CHAIN
    abi *abi.ABI
    topic string
    topicName string
}

func (h *ReferralRelationSubscribe) Init(publisher Publisher) error {
    h.config = publisher.GetConfig()
    h.chain = publisher.GetChain()
    abiBytes, err := os.ReadFile(h.chain.Abi["Poop"])
    if err != nil {
        return err
    }
    abiStr := string(abiBytes)
    
    abi, err := abi.JSON(strings.NewReader(abiStr))
    if err != nil {
        return err
    }
    h.abi = &abi
    
    h.topicName = "ReferralRelation"
    topicSig := []byte("ReferralRelation(address,address,uint256)")
    h.topic = crypto.Keccak256Hash(topicSig).Hex()

    err = publisher.Subscribe(h, h.topic)
    if err != nil {
        return err
    }
    return nil
}

type ReferralRelationEvent struct {
    User        common.Address
    Referral    common.Address
    Time        *big.Int
}

func (h *ReferralRelationSubscribe) Handle(event *types.Log, dbTransaction *gorm.DB) (error) {
    log.Debugf("[%s] handle ReferralRelation event", h.chain.ChainName)
    referralRelationEvent := ReferralRelationEvent{}
    err := h.abi.UnpackIntoInterface(&referralRelationEvent, "ReferralRelation", event.Data)
    if err != nil {
        return err
    }
    log.Debugf("referralRelationEvent: %v", referralRelationEvent)
    //首先更新被邀请者的信息
    userChainInfo, err := UserChainInfo.GetByAddress(dbTransaction, h.chain.ChainName, referralRelationEvent.User.Hex())
    if err != nil {
        return err
    }
    userChainInfoExists := !(userChainInfo == nil)
    if !userChainInfoExists {
        userChainInfo = &UserChainInfo.UserChainInfo {
            Chain: h.chain.ChainName,
            Address: referralRelationEvent.User.Hex(),
            Referral: utils.PTR(referralRelationEvent.Referral.Hex()),
            ReferralTime: utils.PTR(time.Unix(referralRelationEvent.Time.Int64(), 0)),
        }
        affectRow, err := UserChainInfo.Save(dbTransaction, userChainInfo)
        if err != nil {
            return err
        } else if affectRow != 1 {
            return errors.New(fmt.Sprintf("insert userChainInfo failed affectRow=%d", affectRow))
        }
    } else {
        userChainInfo.Referral = utils.PTR(referralRelationEvent.Referral.Hex())
        userChainInfo.ReferralTime = utils.PTR(time.Unix(referralRelationEvent.Time.Int64(), 0))
        affectRow, err := UserChainInfo.UpdateReferralInfo(dbTransaction, userChainInfo)
        if err != nil {
            return err
        } else if affectRow != 1 {
            return errors.New(fmt.Sprintf("update userChainInfo failed affectRow=%d", affectRow))
        }
    }
    //接下来更新邀请者的邀请数量
    userChainInfo, err = UserChainInfo.GetByAddress(dbTransaction, h.chain.ChainName, referralRelationEvent.Referral.Hex())
    if err != nil {
        return err
    }
    if userChainInfo == nil {
        userChainInfo = &UserChainInfo.UserChainInfo {
            Chain: h.chain.ChainName,
            Address: referralRelationEvent.Referral.Hex(),
            InviteNum: uint(1),
        }
        affectRow, err := UserChainInfo.Save(dbTransaction, userChainInfo)
        if err != nil {
            return err
        } else if affectRow != 1 {
            return errors.New(fmt.Sprintf("insert userChainInfo failed affectRow=%d", affectRow))
        }
    } else {
        userChainInfo.InviteNum = userChainInfo.InviteNum + 1
        affectRow, err := UserChainInfo.UpdateInviteNum(dbTransaction, userChainInfo)
        if err != nil {
            return err
        } else if affectRow != 1 {
            return errors.New(fmt.Sprintf("update userChainInfo failed affectRow=%d", affectRow))
        }
    }

    return nil
}
