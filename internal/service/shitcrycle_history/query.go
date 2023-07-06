package shitcrycle_history

import (
    "poop.fi/poop-server/internal/database"
)

func CountByAddress(address string) (int64, error) {
    var num int64
    err := database.GetDB().
        Model(&ShitcrycleHistory{}).
        Where("`address` = ?", address).
        Count(&num).
        Error;
    if err == nil {
        return num, nil
    } else {
        return 0, err
    }
}

