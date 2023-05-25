package server

import (
    log "github.com/sirupsen/logrus"
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "poop.fi/poop-server/internal/config"
    "poop.fi/poop-server/internal/server/routers"
)

func Run(cfg * config.Config) {
    gin.SetMode(cfg.HTTP.Mode)
    
    routersInit := routers.InitRouter()

    endPoint := fmt.Sprintf(":%d", cfg.HTTP.Port)
    maxHeaderBytes := 1 << 20
    server := &http.Server{
        Addr:               endPoint,
        Handler:            routersInit,
        ReadTimeout:        cfg.HTTP.ReadTimeout,
        WriteTimeout:       cfg.HTTP.WriteTimeout,
        MaxHeaderBytes:     maxHeaderBytes,
    }

    log.Infof("start http server at: %s", endPoint)
    server.ListenAndServe()
    
}
