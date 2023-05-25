package database

import (
    //log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "fmt"
)

var db *gorm.DB

func InitDatabase(cfg *config.Config) (error) {
    url := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", 
        cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database, cfg.DB.Timeout,
    )
    var err error
    db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
    if err != nil {
        return err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)

    return nil
}

func GetDB() *gorm.DB {
    return db
}
