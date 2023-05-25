package mlog

import (
    log "github.com/sirupsen/logrus"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "poop.fi/poop-server/internal/config"
    "io"
    //"io/ioutil"
    "os"
    "time"
)


func InitLog(cfg *config.Config) (error) {
    log.SetLevel(cfg.LOG.Level)

    writers := []io.Writer{}

    if cfg.LOG.Console {
        writers = append(writers, os.Stdout)
    }

    if cfg.LOG.File != "" {
        writer, err := rotatelogs.New(
            cfg.LOG.File + ".%Y%m%d",
            rotatelogs.WithLinkName(cfg.LOG.File),
            rotatelogs.WithRotationCount(cfg.LOG.Keep),
            rotatelogs.WithRotationTime(time.Hour * 24),
        )
        if err != nil {
            return err
        }
        writers = append(writers, writer)
    }
    log.SetOutput(io.MultiWriter(writers...))
    return nil
}
