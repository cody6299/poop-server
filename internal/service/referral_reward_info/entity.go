package referral_reward_info

import (
    "time"
    "math/big"
)

type ReferralRewardInfo struct {
    Id              uint64      `gorm:"primaryKey"`
    Chain           string
    Address         string
    InviteAddress   string
    RewardAmount    uint64
    RewardTime      time.Time
    TxHash          string
    CreateAt        time.Time   `gorm:"autoCreateTime"`
    UpdateAt        time.Time   `gorm:"autoUpdateTime"`
}

type AggregationRecord struct {
    Address             string
    TotalRewardAmount   *big.Int
    TotalRewardNum      uint64
}

func (v *ReferralRewardInfo) TableName() string {
    return "referral_reward_info"
}

func (v *AggregationRecord) TableName() string {
    return "referral_reward_info"
}
