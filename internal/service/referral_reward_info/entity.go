package referral_reward_info

import (
    "time"
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

func (v *ReferralRewardInfo) TableName() string {
    return "referral_reward_info"
}
