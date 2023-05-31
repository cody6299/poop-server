package price_info

import (
    "gorm.io/gorm"
)

func SaveOrUpdate(dbTransaction *gorm.DB, priceInfo *PriceInfo) (int64, error) {
    query := "INSERT INTO price_info(`chain`, `price_type`, `price_key`, `begin_time`, `end_time`, `price_open`, `price_high`, `price_low`, `price_close`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `price_high` = VALUES(`price_high`), `price_low` = VALUES(`price_low`), `price_close` = VALUES(`price_close`)"
    values := []interface{} {priceInfo.Chain, priceInfo.PriceType, priceInfo.PriceKey, priceInfo.BeginTime, priceInfo.EndTime, priceInfo.PriceOpen, priceInfo.PriceHigh, priceInfo.PriceLow, priceInfo.PriceClose}
    result := dbTransaction.Exec(query, values...)

    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}
