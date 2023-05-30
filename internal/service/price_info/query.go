package price_info

import (
    "errors"
    "gorm.io/gorm"
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
