package subscribers

import (
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/config"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type Subscriber interface {
    Handle(event *types.Log, dbTransaction *gorm.DB, client *ethclient.Client) error
}

type Publisher interface {
    GetChain() *config.CHAIN
    GetConfig() *config.Config
    Subscribe(subscriber Subscriber, topic string) error
}
