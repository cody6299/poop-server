package process

import (
    "time"
)

type Process struct {
    Id uint64
    Key string
    Value string
    CreateAt time.Time
    UpdateAt time.Time
}

func (v Process) TableName() string {
    return "process"
}
