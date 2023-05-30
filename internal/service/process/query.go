package process

import (
    "errors"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetByKey(key string) (*Process, error) {
    process := Process{} 
    err := database.GetDB().
        Where("`key` = ?", key).
        Take(&process).
        Error;
    if err == nil {
        return &process, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}

func GetByKeyAndDBTransaction(dbTransaction *gorm.DB, key string) (*Process, error) {
    process := Process{} 
    err := dbTransaction.
        Where("`key` = ?", key).
        Take(&process).
        Error;
    if err == nil {
        return &process, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}
