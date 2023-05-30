package events

import (
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/config"
)

func Run(cfg *config.Config) (error) {
    for key, value := range cfg.CHAINS {
        log.Debugf("chainName: %s chainId: %d", value.ChainName, value.ChainId)
        publisher, err := NewPublisher(key, cfg)
        if err != nil {
            return err
        }
        publisher.Run()
    }
    select {}
    return nil
}
