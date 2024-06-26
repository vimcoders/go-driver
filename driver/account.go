// 不允许调用标准库外的包，防止循环引用
package driver

import "time"

type Account struct {
	UserId   int64  `gorm:"primarykey"`
	Passport string `gorm:"unique"`
	Pwd      string
	Created  time.Time `gorm:"comment:创建时间"`
}
