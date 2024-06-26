// 不允许调用标准库外的包，防止循环引用
package driver

import "time"

type User struct {
	Id      string    `bson:"_id"`
	UserId  int64     `bson:"user_id,omitempty"`
	Created time.Time `bson:"created,omitempty"`
	Role    `bson:"role,omitempty"`
}

func (x *User) DocumentId() string {
	return x.Id
}

func (x *User) DocumentName() string {
	return "user"
}

type Role struct {
	RoleId int64  `bson:"role_id"`
	Level  int32  `bson:"level"`
	Exp    int32  `bson:"exp"`
	Name   string `bson:"name"`
	Items  `bson:"items"`
	Tasks  `bson:"tasks"`
}
