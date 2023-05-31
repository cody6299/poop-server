package user_referral_code

import (
    "time"
)

type UserReferralCode struct {
    Id              uint64      `gorm:"primaryKey"`
    Address         string
    ReferralCode    *string
    CreateAt        time.Time   `gorm:"autoCreateTime"`
    UpdateAt        time.Time   `gorm:"autoUpdateTime"`
}

func (v *UserReferralCode) TableName() string {
    return "user_referral_code"
}
