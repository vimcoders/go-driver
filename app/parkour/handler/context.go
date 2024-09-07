package handler

import (
	"context"
	"go-driver/app/parkour/driver"
	"go-driver/log"
	"go-driver/mongox"

	"github.com/google/uuid"
)

type Context struct {
	*driver.User
	context.CancelFunc
	*mongox.Mongo
}

func (x *Context) Update() {
	if err := x.Mongo.Update(x.User); err != nil {
		log.Error(err.Error())
	}
}

func (x *Context) Insert() {
	x.User.Id = uuid.NewString()
	if err := x.Mongo.Insert(x.User); err != nil {
		log.Error(err.Error())
	}
}
