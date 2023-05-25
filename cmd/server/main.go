package main

import (
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/mlog"
    "poop.fi/poop-server/internal/config"
    "poop.fi/poop-server/internal/database"
    "poop.fi/poop-server/internal/server"
)

func main() {
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("config error: %s", err)
    }

    err = mlog.InitLog(cfg)
    if err != nil {
        log.Fatalf("log error: %s", err)
    }

    err = database.InitDatabase(cfg)
    if err != nil {
        log.Fatalf("database error: %s", err)
    }

    server.Run(cfg)
}
