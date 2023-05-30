package user_chain_info

import (
    "gorm.io/gorm"
)

func Save(dbTransaction *gorm.DB, record *UserChainInfo) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := dbTransaction.Create(&record)
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}

func UpdateUserId(dbTransaction *gorm.DB, record *UserChainInfo) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := dbTransaction.Model(record).
        Where("id = ?", record.Id).
        Updates(map[string]interface{}{
            "user_id": record.UserId,
        })
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}

func UpdateReferralInfo(dbTransaction *gorm.DB, record *UserChainInfo) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := dbTransaction.Model(record).
        Where("id = ?", record.Id).
        Updates(map[string]interface{}{
            "referral": record.Referral,
            "referral_time": record.ReferralTime,
        })
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}

func UpdateInviteNum(dbTransaction *gorm.DB, record *UserChainInfo) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := dbTransaction.Model(record).
        Where("id = ?", record.Id).
        Updates(map[string]interface{}{
            "invite_num": record.InviteNum,
        })
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}

func UpdateRewardInfo(dbTransaction *gorm.DB, record *UserChainInfo) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := dbTransaction.Model(record).
        Where("id = ?", record.Id).
        Updates(map[string]interface{}{
            "invite_reward": record.InviteReward,
            "reward_num": record.RewardNum,
        })
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}
