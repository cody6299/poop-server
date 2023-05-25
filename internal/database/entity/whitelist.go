package entity

import (
    "time"
    "poop.fi/poop-server/internal/database"
    "gorm.io/gorm"
    "errors"
)

type Whitelist struct {
    Id uint64
    Address string
    MaxAmount string
    Proof string
    CreateAt time.Time
    UpdateAt time.Time
}

func (v Whitelist) TableName() string {
    return "whitelist"
}

func GetByAddress(address string) (*Whitelist, error) {
    whitelist := Whitelist{} 
    err := database.GetDB().
        Where("address = ?", address).
        Take(&whitelist).
        Error;
    if err == nil {
        return &whitelist, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}
