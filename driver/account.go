package driver

import "time"

type Account struct {
	UserId   int64  `gorm:"primarykey"`
	Passport string `gorm:"unique"`
	Pwd      string
	Created  time.Time `gorm:"comment:创建时间"`
}
