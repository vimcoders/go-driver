package handler

import (
	"github.com/vimcoders/go-driver/app/parkour/driver"

	"github.com/vimcoders/go-driver/mongox"

	"github.com/vimcoders/go-driver/log"

	"go.mongodb.org/mongo-driver/bson"
)

func (x *Handler) GetUser(userId int64) *driver.User {
	x.Lock()
	defer x.Unlock()
	for i := 0; i < len(x.Users); i++ {
		if x.Users[i].UserId == userId {
			return x.Users[i]
		}
	}
	users, err := mongox.Query[*driver.User](x.Mongo, bson.M{"user_id": userId})
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	x.Users = append(x.Users, users...)
	for i := 0; i < len(users); i++ {
		if users[i].UserId == userId {
			return users[i]
		}
	}
	return nil
}
