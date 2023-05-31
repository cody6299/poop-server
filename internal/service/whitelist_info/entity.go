package whitelist_info

import (
    "time"
)

type WhitelistInfo struct {
    Id          uint64      `gorm:"primaryKey"`
    Chain       string
    Address     string
    MaxAmount   string
    Proof       string
    CreateAt    time.Time   `gorm:"autoCreateTime"`
    UpdateAt    time.Time   `gorm:"autoUpdateTime"`
}

func (v WhitelistInfo) TableName() string {
    return "whitelist_info"
}
