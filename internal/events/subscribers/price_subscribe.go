package subscribers

import (
    "os"
    "fmt"
    "time"
    "errors"
    "strconv"
    "strings"
    "math/big"
    "gorm.io/gorm"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/config"
    "poop.fi/poop-server/internal/service/process"
    "poop.fi/poop-server/internal/utils"
    UserChainInfo "poop.fi/poop-server/internal/service/user_chain_info"
    PriceInfo "poop.fi/poop-server/internal/service/price_info"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/accounts/abi"
)

type PriceSubscribe struct {
    config *config.Config
    chain *config.CHAIN
    abi *abi.ABI
    topic string
    topicName string
}

func (h *PriceSubscribe) Init(publisher Publisher) error {
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
    
    h.topicName = "Price"
    topicSig := []byte("Price(address,uint256,uint256,uint256)")
    h.topic = crypto.Keccak256Hash(topicSig).Hex()

    err = publisher.Subscribe(h, h.topic)
    if err != nil {
        return err
    }
    return nil
}

type PriceEvent struct {
    User        common.Address
    Time        *big.Int
    Recieved    *big.Int
    Sent        *big.Int
}

func (h *PriceSubscribe) Handle(event *types.Log, dbTransaction *gorm.DB) (error) {
    log.Debugf("[%s] handle Price event", h.chain.ChainName)
    priceEvent := PriceEvent{}
    err := h.abi.UnpackIntoInterface(&priceEvent, "Price", event.Data)
    if err != nil {
        return err
    }
    log.Debugf("priceEvent: %v", priceEvent)
    //首先判断用户是否存在,如果是第一次交易为其生成uid
    userChainInfo, err := UserChainInfo.GetByAddress(dbTransaction, h.chain.ChainName, priceEvent.User.Hex())
    if err != nil {
        return err
    }
    userChainInfoExists := !(userChainInfo == nil)
    if !userChainInfoExists || userChainInfo.UserId == nil {
        currentMaxUserId, err := process.GetByKeyAndDBTransaction(dbTransaction, fmt.Sprintf("%s_userid", h.chain.ChainName))
        if err != nil {
            return err
        }
        if currentMaxUserId == nil {
            currentMaxUserId = &process.Process{
                Key: fmt.Sprintf("%s_userid", h.chain.ChainName),
                Value: fmt.Sprintf("%d", 0),
            }
        }
        nextUserId, err := strconv.ParseUint(currentMaxUserId.Value, 10, 64)
        if err != nil {
            return err
        }
        nextUserId = nextUserId + 1

        if (!userChainInfoExists) {
            userChainInfo = &UserChainInfo.UserChainInfo{
                Chain: h.chain.ChainName,
                Address: priceEvent.User.Hex(),
                UserId: utils.PTR(nextUserId),
            }
            affectRow, err := UserChainInfo.Save(dbTransaction, userChainInfo)
            if err != nil {
                return err
            } else if affectRow != 1 {
                return errors.New(fmt.Sprintf("save userChainInfo failed affectRow=%d", affectRow)) 
            }
        } else if userChainInfo.UserId == nil {
            userChainInfo.UserId = utils.PTR(nextUserId)
            affectRow, err := UserChainInfo.UpdateUserId(dbTransaction, userChainInfo)
            if err != nil {
                return err
            } else if affectRow != 1 {
                return errors.New(fmt.Sprintf("update userChainInfo failed affectRow=%d", affectRow)) 
            }
        }
        currentMaxUserId.Value = fmt.Sprintf("%d", nextUserId)
        affectRow, err := process.InsertOrUpdate(dbTransaction, currentMaxUserId)
        if err != nil {
            return err
        } else if affectRow != 1 && affectRow != 2 {
            return errors.New(fmt.Sprintf("insert or update process failed affectRow=%d", affectRow)) 
        }
        log.Infof("[%s] update user [address=%s] [user_id=%d]",  h.chain.ChainName, priceEvent.User.Hex(), nextUserId)
    }
    //然后更新价格维度
    priceTime := time.Unix(priceEvent.Time.Int64(), 0).UTC()
    log.Debugf("priceTime: %v", priceTime)
    tmp := big.NewInt(0)
    tmp.Mul(priceEvent.Sent, big.NewInt(1000000000000000000))
    price := big.NewInt(0)
    price.Div(tmp, priceEvent.Recieved)
    log.Debugf("price: %s", price.String())
    {
        //1minute
        affectRow, err := h.priceHistory(int64(60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //5minute
        affectRow, err := h.priceHistory(int64(5 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //15minute
        affectRow, err := h.priceHistory(int64(15 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //30minute
        affectRow, err := h.priceHistory(int64(30 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //1 hour
        affectRow, err := h.priceHistory(int64(60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //4 hour
        affectRow, err := h.priceHistory(int64(4 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //1 day
        affectRow, err := h.priceHistory(int64(24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //3 day
        affectRow, err := h.priceHistory(int64(3 * 24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //7 day
        affectRow, err := h.priceHistory(int64(7 * 24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //14 day
        affectRow, err := h.priceHistory(int64(14 * 24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //1 month
        affectRow, err := h.priceHistory(int64(30 * 24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //3 month
        affectRow, err := h.priceHistory(int64(3 * 30 * 24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }
    {
        //6 month
        affectRow, err := h.priceHistory(int64(6 * 30 * 24 * 60 * 60), priceTime, price.Uint64(), dbTransaction)
        if err != nil {
            return err
        } else {
            log.Debugf("process 1minute price history finish. affectRow=%d", affectRow)
        }
    }

    return nil
}

func (h *PriceSubscribe) priceHistory(interval int64, priceTime time.Time, price uint64, dbTransaction *gorm.DB) (int64, error) {
    priceType := fmt.Sprintf("price_%d", interval)
    keyTime := uint64(priceTime.Add(-time.Duration(priceTime.Unix() % interval) * time.Second).Unix()) 
    log.Debugf("priceHistory1Minute keyTime: %d", keyTime)
    priceInfo, err := PriceInfo.GetByChainAndTypeAndKey(dbTransaction, h.chain.ChainName, priceType, keyTime)
    if err != nil {
        return 0, nil
    }
    if priceInfo == nil {
        priceInfo = &PriceInfo.PriceInfo{
            Chain: h.chain.ChainName,
            PriceType: priceType,
            PriceKey: keyTime,
            BeginTime: keyTime,
            EndTime: keyTime + uint64(interval - 1),
            PriceOpen: price,
            PriceHigh: price,
            PriceLow: price,
            PriceClose: price,
        }
    } else {
        if price > priceInfo.PriceHigh {
            priceInfo.PriceHigh = price
        }
        if price < priceInfo.PriceLow {
            priceInfo.PriceLow = price
        }
        priceInfo.PriceClose = price
    }
    affectRow, err := PriceInfo.SaveOrUpdate(dbTransaction, priceInfo)
    if err != nil {
        return 0, err
    } else if affectRow != 1 && affectRow != 2 {
        //return 0, errors.New(fmt.Sprintf("insert or update priceinfo failed affectRow=%d", affectRow))
    }
    return affectRow, nil
}
