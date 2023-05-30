package user_chain_info

import (
    "errors"
    "gorm.io/gorm"
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
