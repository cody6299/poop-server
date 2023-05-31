package user_referral_code

import (
    "poop.fi/poop-server/internal/database"
)

func Save(record *UserReferralCode) (int64, error) {
    if record == nil {
        return 0, nil
    }
    result := database.GetDB().Create(&record)
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}

func UpdateReferralCode(record *UserReferralCode) (int64, error) {
    result := database.GetDB().Model(record).
        Where("id = ?", record.Id).
        Updates(map[string]interface{}{
            "referral_code": record.ReferralCode,
        })
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}
