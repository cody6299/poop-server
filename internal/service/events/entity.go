package events

import (
    "time"
)

type Events struct {
    Id          uint64      `gorm:"primaryKey"`
    Chain       string
    BlockHash   string
    TxIndex     uint
    LogIndex    uint
    Contract    string
    Topics      string
    Data        string
    BlockNumber uint
    TxHash      string
    Removed     bool
    CreateAt    time.Time   `gorm:"autoCreateTime"`
    UpdateAt    time.Time   `gorm:"autoUpdateTime"`
}

func (v Events) TableName() string {
    return "events"
}
