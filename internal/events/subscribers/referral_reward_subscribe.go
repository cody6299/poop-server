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
    ReferralRewardInfo "poop.fi/poop-server/internal/service/referral_reward_info"
    UserChainInfo "poop.fi/poop-server/internal/service/user_chain_info"
)

type ReferralRewardSubscribe struct {
    config *config.Config
    chain *config.CHAIN
    abi *abi.ABI
    topic string
    topicName string
}

func (h *ReferralRewardSubscribe) Init(publisher Publisher) error {
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
    
    h.topicName = "ReferralReward"
    topicSig := []byte("ReferralReward(address,address,uint256,uint256)")
    h.topic = crypto.Keccak256Hash(topicSig).Hex()

    err = publisher.Subscribe(h, h.topic)
    if err != nil {
        return err
    }
    return nil
}

type ReferralRewardEvent struct {
    User        common.Address
    Referral    common.Address
    Amount      *big.Int
    Time        *big.Int
}

func (h *ReferralRewardSubscribe) Handle(event *types.Log, dbTransaction *gorm.DB) (error) {
    log.Debugf("[%s] handle ReferralReward event", h.chain.ChainName)
    referralRewardEvent := ReferralRewardEvent{}
    err := h.abi.UnpackIntoInterface(&referralRewardEvent, "ReferralReward", event.Data)
    if err != nil {
        return err
    }
    log.Debugf("referralRewardEvent: %v", referralRewardEvent)
    //首先存储记录
    referralRewardInfo := &ReferralRewardInfo.ReferralRewardInfo{
        Chain: h.chain.ChainName,
        Address: referralRewardEvent.Referral.Hex(),
        InviteAddress: referralRewardEvent.User.Hex(),
        RewardAmount: referralRewardEvent.Amount.Uint64(),
        RewardTime: time.Unix(referralRewardEvent.Time.Int64(), 0),
        TxHash: event.TxHash.Hex(),
    }
    affectRow, err := ReferralRewardInfo.Save(dbTransaction, referralRewardInfo)
    if err != nil {
        return err
    }
    if affectRow != 1 {
        return errors.New(fmt.Sprintf("save ReferralRewardInfo failed affectRow=%d", affectRow))
    }
    //接下来更新UserChainInfo
    userChainInfo, err := UserChainInfo.GetByAddress(dbTransaction, h.chain.ChainName, referralRewardEvent.Referral.Hex())
    if err != nil {
        return err
    }
    if userChainInfo == nil {
        userChainInfo = &UserChainInfo.UserChainInfo{
            Chain: h.chain.ChainName,
            Address: referralRewardEvent.Referral.Hex(),
            InviteReward: referralRewardEvent.Amount.Uint64(),
            RewardNum: 1,
        }
        affectRow, err = UserChainInfo.Save(dbTransaction, userChainInfo)
        if err != nil {
            return err
        }
        if affectRow != 1 {
            return errors.New(fmt.Sprintf("save UserChainInfo failed affectRow=%d", affectRow))
        }
    } else {
        userChainInfo.InviteReward = userChainInfo.InviteReward + referralRewardEvent.Amount.Uint64()
        userChainInfo.RewardNum = userChainInfo.RewardNum + 1
        affectRow, err = UserChainInfo.UpdateRewardInfo(dbTransaction, userChainInfo)
        if err != nil {
            return err
        }
        if affectRow != 1 {
            return errors.New(fmt.Sprintf("update UserChainInfo failed affectRow=%d", affectRow))
        }
    }
    return nil
}

