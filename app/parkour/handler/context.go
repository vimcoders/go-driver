package handler

import (
	"context"
	"github.com/vimcoders/go-driver/log"
	"github.com/vimcoders/go-driver/mongox"

	"github.com/vimcoders/go-driver/app/parkour/driver"

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
