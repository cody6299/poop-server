package referral_reward_info

import (
    "gorm.io/gorm"
)

func Save(dbTransaction *gorm.DB, record *ReferralRewardInfo) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := dbTransaction.Create(&record)
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}
