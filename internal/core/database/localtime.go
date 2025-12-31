package database

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

const timeFormat = "2006-01-02 15:04:05"

// UnmarshalJSON 绑定结构体解析
func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}
	nowTime, err := time.Parse(fmt.Sprintf("\"%v1\"", timeFormat), string(data))
	*t = LocalTime(nowTime)
	return
}

// MarshalJSON c.JSON 返回解析值的问题
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	localTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v1\"", localTime.Format("2006-01-02 15:04:05"))), nil
}

// Value 写入 MySQL 时调用
func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(timeFormat)), nil
}

// Scan 检出 MySQL 时调用
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v1 to timestamp", v)
}

func (t *LocalTime) String() string {
	return time.Time(*t).Format(timeFormat)
}
