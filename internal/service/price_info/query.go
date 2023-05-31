package price_info

import (
    "errors"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetByChainAndTypeAndKey(dbTransaction *gorm.DB, chain string, priceType string, priceKey uint64) (*PriceInfo, error) {
    priceInfo := PriceInfo{} 
    err := dbTransaction.
        Where("`chain` = ?", chain).
        Where("`price_type` = ?", priceType).
        Where("`price_key` = ?", priceKey).
        Take(&priceInfo).
        Error;
    if err == nil {
        return &priceInfo, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}

func GetRangeByChainAndType(chain string, priceType string, offset uint, limit uint) (*[]PriceInfo, error) {
    var records []PriceInfo
    err := database.GetDB().
        Where("`chain` = ?", chain).
        Where("`price_type` = ?", priceType).
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

func GetByChainAndType(chain string, priceType string) (*[]PriceInfo, error) {
    var records []PriceInfo
    err := database.GetDB().
        Where("`chain` = ?", chain).
        Where("`price_type` = ?", priceType).
        Order("id desc").
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
