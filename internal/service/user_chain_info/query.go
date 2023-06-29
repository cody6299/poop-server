package user_chain_info

import (
    "time"
    "errors"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetByAddress(dbTransaction *gorm.DB, chain string, address string) (*UserChainInfo, error) {
    userChainInfo := UserChainInfo{} 
    err := dbTransaction.
        Where("`address` = ?", address).
        Where("`chain` = ?", chain).
        Take(&userChainInfo).
        Error;
    if err == nil {
        return &userChainInfo, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}

func GetByChainAndAddress(chain string, address string) (*UserChainInfo, error) {
    userChainInfo := UserChainInfo{} 
    err := database.GetDB().
        Where("`address` = ?", address).
        Where("`chain` = ?", chain).
        Take(&userChainInfo).
        Error;
    if err == nil {
        return &userChainInfo, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}

func AggregationByTime(chain string, startTime int, endTime int) (*[]AggregationRecord, error) {
    var records []AggregationRecord

    err := database.GetDB().
        Table("user_chain_info").
        Select("referral, count(1) as total_referral_num").
        Where("`chain` = ?", chain).
        Where("`referral_time` >= ?", time.Unix(int64(startTime), 0)).
        Where("`referral_time` <= ?", time.Unix(int64(endTime), 0)).
        Group("referral").
        Scan(&records).
        Error
    if err == nil {
        return &records, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return &records, nil
    } else {
        return nil, err
    }
}
