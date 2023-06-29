package referral_reward_info

import (
    "time"
    "errors"
    "math/big"
    "gorm.io/gorm"
    "poop.fi/poop-server/internal/database"
)

func GetRangeByChainAndAddress(chain string, address string, offset uint, limit uint) (*[]ReferralRewardInfo, error) {
    var records []ReferralRewardInfo
    err := database.GetDB().
        Where("`chain` = ?", chain).
        Where("`address` = ?", address).
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

func AggregationByTime(chain string, startTime int, endTime int) (*[]AggregationRecord, error) {
    var records []ReferralRewardInfo
    err := database.GetDB().
        Table("referral_reward_info").
        Select("address, reward_amount").
        Where("`chain` = ?", chain).
        Where("`reward_time` >= ?", time.Unix(int64(startTime), 0)).
        Where("`reward_time` <= ?", time.Unix(int64(endTime), 0)).
        Scan(&records).
        Error
    if err != nil {
        return nil, err
    }
    var recordAmountMap map[string]*big.Int = make(map[string]*big.Int)
    var recordNumMap map[string]uint64 = make(map[string]uint64)
    var addresses []string
    for _, record := range records {
        address := record.Address
        if recordAmountMap[address] == nil {
            recordAmountMap[address] = new(big.Int).SetUint64(0)
            addresses = append(addresses, address)
        }
        recordAmountMap[address].Add(recordAmountMap[address], new (big.Int).SetUint64(record.RewardAmount))
        recordNumMap[address] = recordNumMap[address] + 1
    }
    var res []AggregationRecord
    for _, address := range addresses {
        res = append(res, AggregationRecord{
            address,
            recordAmountMap[address],
            recordNumMap[address],
        })
    }
    return &res, nil
}
