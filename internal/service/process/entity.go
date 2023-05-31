package process

import (
    "time"
)

type Process struct {
    Id          uint64      `gorm:"primaryKey"`
    Key         string
    Value       string
    CreateAt    time.Time   `gorm:"autoCreateTime"`
    UpdateAt    time.Time   `gorm:"autoUpdateTime"`
}

func (v Process) TableName() string {
    return "process"
}
