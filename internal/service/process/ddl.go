package process

import (
    "gorm.io/gorm"
)

func InsertOrUpdate(dbTransaction *gorm.DB, process *Process) (int64, error) {
    query := "INSERT INTO process(`key`, `value`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = VALUES(`value`)"
    values := []interface{} {process.Key, process.Value}
    result := dbTransaction.Exec(query, values...)

    if result.Error != nil {
        return 0, result.Error
    }
    return result.RowsAffected, nil
}
