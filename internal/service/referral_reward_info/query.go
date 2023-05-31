package referral_reward_info

import (
    "errors"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetRangeByChainAndAddress(chain string, address string, offset uint, limit uint) (*[]ReferralRewardInfo, error) {
    var records []ReferralRewardInfo
    err := database.GetDB().
        Where("`chain` = ?", chain).
        Where("`address` = ?", address).
        Order("id desc").
        Offset(int(offset)).
        Limit(int(limit)).
        Find(&records).
        Error;
    if err == nil {
        return &records, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return &records, nil
    } else {
        return nil, err
    }
}
