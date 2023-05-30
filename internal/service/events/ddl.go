package events

import (
    "gorm.io/gorm"
)

func SaveBulk(dbTransaction *gorm.DB, records []*Events) (int64, error) {
    if len(records) == 0 {
        return 0, nil
    }
    result := dbTransaction.Create(&records)
    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}
