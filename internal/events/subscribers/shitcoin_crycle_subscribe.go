package subscribers

import (
    "os"
    "fmt"
    "context"
    "time"
    "errors"
    //"strconv"
    "strings"
    "math/big"
    "gorm.io/gorm"
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/config"
    //"poop.fi/poop-server/internal/service/process"
    //"poop.fi/poop-server/internal/utils"
    //UserChainInfo "poop.fi/poop-server/internal/service/user_chain_info"
    //PriceInfo "poop.fi/poop-server/internal/service/price_info"
    ShitcrycleHistory "poop.fi/poop-server/internal/service/shitcrycle_history"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/ethclient"
)

type ShitcoinCrycleSubscribe struct {
    config *config.Config
    chain *config.CHAIN
    abi *abi.ABI
    topic string
    topicName string
}

func (h *ShitcoinCrycleSubscribe) Init(publisher Publisher) error {
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
    
    h.topicName = "ShitcoinCrycle"
    topicSig := []byte("Price(address,uint256,uint256,uint256)")
    h.topic = crypto.Keccak256Hash(topicSig).Hex()

    err = publisher.Subscribe(h, h.topic)
    if err != nil {
        return err
    }
    return nil
}

var transferTopic string = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex()
var priceTopic string = crypto.Keccak256Hash([]byte("Price(address,uint256,uint256,uint256)")).Hex()

func (h *ShitcoinCrycleSubscribe) Handle(event *types.Log, dbTransaction *gorm.DB, client *ethclient.Client) (error) {
    log.Debugf("[%s] handle Price event for ShitcoinCrycle record", h.chain.ChainName)
    txHash := event.TxHash
    log.Debugf("transaction: %v", txHash)
    transactionRecipient, err := client.TransactionReceipt(context.Background(), txHash)
    if err != nil {
        return err
    }
    log.Debugf("transactionRecipient: %T", transactionRecipient)
    logs := transactionRecipient.Logs
    //0表示搜索Transfer事件 1表示搜索Price事件
    var eventTypeToSearch = 0;
    var costAmount *big.Int = nil
    var token *common.Address = nil
    shitcrycleHistoryRecords := []*ShitcrycleHistory.ShitcrycleHistory{}
    for _, vLog := range logs {
        topic := vLog.Topics[0].Hex()
        if eventTypeToSearch == 0 {
            if topic != transferTopic {
                continue
            }
            dataType, err := abi.NewType("address", "", nil)
            if err != nil {
                return err
            }
            decodeABI := abi.Arguments{
                {Type: dataType},
            }
            toAddrData, err := decodeABI.Unpack(vLog.Topics[2].Bytes())
            if err != nil {
                return err
            }
            toAddr := fmt.Sprintf("%v", toAddrData[0])
            log.Debugf("to: %T %s %s", toAddr, toAddr, h.chain.Addresses.PoopRouter)
            if toAddr != h.chain.Addresses.PoopRouter && toAddr != h.chain.Addresses.WBNBRouter {
                continue
            }
            token = &vLog.Address
            dataType, err = abi.NewType("uint256", "", nil)
            if err != nil {
                return err
            }
            decodeABI = abi.Arguments{
                {Type: dataType},
            }
            amount, err := decodeABI.Unpack(vLog.Data)
            if err != nil {
                return err
            }
            costAmount, _ = new (big.Int).SetString(fmt.Sprintf("%v", amount[0]), 10)
            log.Debugf("transfer info: %v, %v, %v", token, amount, costAmount)
            eventTypeToSearch = 1
        } else if eventTypeToSearch == 1 {
            if topic != priceTopic {
                continue
            }
            priceEvent := PriceEvent{}
            err := h.abi.UnpackIntoInterface(&priceEvent, "Price", event.Data)
            if err != nil {
                return err
            }
            shitcrycleHistoryRecord := &ShitcrycleHistory.ShitcrycleHistory{
                Chain: h.chain.ChainName,
                Address: priceEvent.User.Hex(),
                Token: token.Hex(),
                Cost: costAmount.String(),
                Recieve: priceEvent.Recieved.String(),
                ActionTime: time.Unix(priceEvent.Time.Int64(), 0).UTC(),
                TxHash: txHash.Hex(),
            }
            shitcrycleHistoryRecords = append(shitcrycleHistoryRecords, shitcrycleHistoryRecord)
            
            token = nil
            costAmount = nil
            eventTypeToSearch = 0
        }
    }
    affectRow, err := ShitcrycleHistory.SaveBulk(dbTransaction, shitcrycleHistoryRecords);
    if err != nil {
        return err;
    }
    if int(affectRow) != len(shitcrycleHistoryRecords) {
        return errors.New(fmt.Sprintf("insert ShitcrycleHistoryRecord failed affectRow=%d dataLength=%d", affectRow, len(shitcrycleHistoryRecords)))
    }
    return nil
}
