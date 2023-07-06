package shitcrycle_history

import (
    "time"
)
type ShitcrycleHistory struct {
    Id          uint64      `gorm:"primaryKey"`
    Chain       string
    Address     string
    Token       string
    Cost        string
    Recieve     string
    ActionTime  time.Time
    TxHash      string
    CreateAt    time.Time   `gorm:"autoCreateTime"`
    UpdateAt    time.Time   `gorm:"autoUpdateTime"`
}

func (v ShitcrycleHistory) TableName() string {
    return "shitcrycle_history"
}
