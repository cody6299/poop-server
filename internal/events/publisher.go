package events

import (
    "fmt"
    "time"
    "strconv"
    "strings"
    "context"
    "math/big"
    "encoding/hex"
    log "github.com/sirupsen/logrus"
    "github.com/ethereum/go-ethereum"
    "poop.fi/poop-server/internal/config"
    "poop.fi/poop-server/internal/database"
    "poop.fi/poop-server/internal/events/subscribers"
    "poop.fi/poop-server/internal/service/process"
    chainEvent "poop.fi/poop-server/internal/service/events"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
)

type Publisher struct {
    config *config.Config
    chain *config.CHAIN
    urls []string
    contracts []common.Address
    subscribers map[string][]subscribers.Subscriber
}

func (publisher *Publisher) GetChain() (*config.CHAIN) {
    return publisher.chain
}

func (publisher *Publisher) GetConfig() (*config.Config) {
    return publisher.config
}

func (publisher *Publisher) Subscribe(subscriber subscribers.Subscriber, topic string) error {
    publisher.subscribers[topic] = append(publisher.subscribers[topic], subscriber)
    return nil
}

func NewPublisher(key string, cfg *config.Config) (*Publisher, error) {
    chain := cfg.CHAINS[key]
    publisher := Publisher{
        subscribers: make(map[string][]subscribers.Subscriber),
        chain: &chain,
        config: cfg,
        urls: []string{},
        contracts: []common.Address{},
    }

    for _, url := range chain.Urls {
        publisher.urls = append(publisher.urls, url)
    }

    for _, contract := range chain.Contracts {
        publisher.contracts = append(publisher.contracts, common.HexToAddress(contract))
    }

    { 
        err := (&subscribers.ReferralRelationSubscribe{}).Init(&publisher)
        if err != nil {
            return nil, err
        }
    }
    {
        err := (&subscribers.ReferralRewardSubscribe{}).Init(&publisher)
        if err != nil {
            return nil, err
        }
    }
    {
        err := (&subscribers.PriceSubscribe{}).Init(&publisher)
        if err != nil {
            return nil, err
        }
    }
    return &publisher, nil
}

func (publisher *Publisher) Run() {
    chainName := publisher.chain.ChainName
    log.Infof("[%s] publisher begin", chainName)
    go func() {
        for {
            log.Infof("[%s] publisher work begin", chainName)

            doSleep := true

            currentProcess, err := process.GetByKey(chainName)
            if err != nil {
                log.Errorf("get current process failed: %v", err)
                break
            }
            if currentProcess == nil {
                currentProcess = &process.Process{
                    Key: chainName,
                    Value: fmt.Sprintf("%d", publisher.chain.StartBlock),
                }
            }

            currentBlockNumber, err := strconv.ParseUint(currentProcess.Value, 10, 64)
            if err != nil {
                log.Errorf("convert process value failed: %v", err)
                break
            }

            availableUrlLength := len(publisher.urls)
            for i := 0; i < availableUrlLength; i ++ {
                url := publisher.urls[0]
                publisher.urls = publisher.urls[1:]
                publisher.urls = append(publisher.urls, url)

                client, err := ethclient.Dial(url)
                if err != nil {
                    continue
                }

                chainBlockNumber, err := client.BlockNumber(context.Background())
                if err != nil {
                    continue
                }

                //chainBlockNumber = 30194453

                log.Debugf("[%s] chain block number: %d", publisher.chain.ChainName, chainBlockNumber)

                endBlockNumber := chainBlockNumber - publisher.chain.DelayBlock
                if (endBlockNumber > currentBlockNumber + 1 + publisher.chain.MaxBlock) {
                    endBlockNumber = currentBlockNumber + 1 + publisher.chain.MaxBlock
                }

                if currentBlockNumber >= endBlockNumber {
                    log.Infof("[%s] process do nothing [currentBlockNumber=%d] [chainBlockNumber=%d] [endBlockNumber=%d]", 
                        chainName, currentBlockNumber, chainBlockNumber, endBlockNumber)
                    break
                }

                query := ethereum.FilterQuery{
                    FromBlock: new(big.Int).SetUint64(currentBlockNumber + 1),
                    ToBlock: new(big.Int).SetUint64(endBlockNumber),
                    Addresses: publisher.contracts,
                }
                logs, err := client.FilterLogs(context.Background(), query)
                if err != nil {
                    continue
                }
                log.Debugf("[%s] query Logs from chain [%d-%d] logs num: %d", chainName, currentBlockNumber + 1, endBlockNumber, len(logs))
                
                //开始事务
                dbTransaction := database.GetDB().Begin()
                //首先存储所有的events
                eventRecords := []*chainEvent.Events{}
                for _, vLog := range logs {
                    vLogTopics := []string{}
                    for _, topic := range(vLog.Topics) {
                        vLogTopics = append(vLogTopics, topic.Hex())
                    }
                    eventRecord := &chainEvent.Events{
                        Chain: chainName,
                        BlockHash: vLog.BlockHash.Hex(),
                        TxIndex: vLog.TxIndex,
                        LogIndex: vLog.Index,
                        Contract: vLog.Address.Hex(),
                        Topics: strings.Join(vLogTopics, ","),
                        Data: hex.EncodeToString(vLog.Data),
                        BlockNumber: uint(vLog.BlockNumber),
                        TxHash: vLog.TxHash.Hex(),
                        Removed: vLog.Removed,
                    }
                    eventRecords = append(eventRecords, eventRecord)
                }
                affectRow, err := chainEvent.SaveBulk(dbTransaction, eventRecords)
                if err != nil || affectRow != int64(len(logs)) {
                    log.Errorf("[%s] database safe events failed, err=%v, affectRow=%d", chainName, err, affectRow)
                    dbTransaction.Rollback()
                    break
                } else {
                    log.Debugf("[%s] save events finish. num=%d", chainName, affectRow)
                }
                //接下来将event交给所有的subscriber
                err = nil
                for _, vLog := range logs {
                    if err != nil {
                        break
                    }
                    topic := vLog.Topics[0].Hex()
                    subscribers, exists := publisher.subscribers[topic]
                    if exists && subscribers != nil {
                        for _, subscriber := range(subscribers) {
                            err = subscriber.Handle(&vLog, dbTransaction)
                            if err != nil {
                                break
                            }
                        }
                    }
                }
                if err != nil {
                    log.Errorf("[%s] handle event failed: %v", chainName, err)
                    dbTransaction.Rollback()
                    break
                }

                currentProcess.Value = fmt.Sprintf("%d", endBlockNumber)
                affectRow, err = process.InsertOrUpdate(dbTransaction, currentProcess)
                if err != nil || (affectRow != 1 && affectRow != 2) {
                    log.Errorf("[%s] database update process failed, err=%v, affectRow=%d", chainName, err, affectRow)
                    dbTransaction.Rollback()
                    break
                } else {
                    log.Debugf("[%s] update process finish. value=%s", chainName, currentProcess.Value)
                }
                //提交事务
                dbTransaction.Commit()
                if endBlockNumber < chainBlockNumber - publisher.chain.DelayBlock {
                    doSleep = false
                }
                log.Infof("[%s] publisher work finish. process block [%d-%d] events: %d, doSleep=%t", chainName, currentBlockNumber + 1, endBlockNumber, len(logs), doSleep)

                break
            }

            if (doSleep) {
                log.Infof("[%s] reach to bestBlockNumber sleep %v", chainName, publisher.chain.Interval)
                time.Sleep(publisher.chain.Interval)
            } else {
                log.Debugf("[%s] more work to do continue request", chainName)
            }
            //break
        }
    }()
}
