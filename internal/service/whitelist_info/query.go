package whitelist_info

import (
    "errors"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetByChainAndAddress(chain string, address string) (*WhitelistInfo, error) {
    whitelistInfo := WhitelistInfo{} 
    err := database.GetDB().
        Where("`chain` = ?", chain).
        Where("`address` = ?", address).
        Take(&whitelistInfo).
        Error;
    if err == nil {
        return &whitelistInfo, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}
