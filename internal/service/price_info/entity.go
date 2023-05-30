package price_info

import (
    "time"
)

type PriceInfo struct {
    Id              uint64      `gorm:"primaryKey"`
    Chain           string
    PriceType       string
    PriceKey        uint64
    PriceOpen       uint64
    PriceHigh       uint64
    PriceLow        uint64
    PriceClose      uint64
    CreateAt        time.Time   `gorm:"autoCreateTime"`
    UpdateAt        time.Time   `gorm:"autoUpdateTime"`
}

func (v PriceInfo) TableName() string {
    return "price_info"
}
