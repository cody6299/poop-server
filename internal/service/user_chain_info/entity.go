package user_chain_info

import (
    "time"
)

type UserChainInfo struct {
    Id              uint64      `gorm:"primaryKey"`
    Chain           string
    Address         string
    UserId          *uint64
    Referral        *string
    ReferralTime    *time.Time
    InviteNum       uint
    InviteReward    uint64
    RewardNum       uint
    CreateAt        time.Time   `gorm:"autoCreateTime"`
    UpdateAt        time.Time   `gorm:"autoUpdateTime"`
}

func (v UserChainInfo) TableName() string {
    return "user_chain_info"
}
