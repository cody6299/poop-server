package user_referral_code

import (
    "errors"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetByAddress(address string) (*UserReferralCode, error) {
    userReferralCode := UserReferralCode{} 
    err := database.GetDB().
        Where("`address` = ?", address).
        Take(&userReferralCode).
        Error;
    if err == nil {
        return &userReferralCode, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}

func GetByCode(code string) (*UserReferralCode, error) {
    userReferralCode := UserReferralCode{} 
    err := database.GetDB().
        Where("`referral_code` = ?", code).
        Take(&userReferralCode).
        Error;
    if err == nil {
        return &userReferralCode, nil
    } else if (errors.Is(err, gorm.ErrRecordNotFound)) {
        return nil, nil
    } else {
        return nil, err
    }
}
