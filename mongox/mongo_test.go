package mongox_test

import (
	"fmt"
	"testing"

	"github.com/vimcoders/go-driver/mongox"
)

type Document struct {
	Id  string
	Say string
}

func (x *Document) DocumentName() string {
	return "document"
}

func (x *Document) DocumentId() string {
	return x.Id
}

func TestUpload(t *testing.T) {
	mongo, err := mongox.Connect("mongodb://admin:admin@127.0.0.1:27017", "parkour")
	if err != nil {
		fmt.Println(err)
	}
	document := Document{
		Id:  "2084226a-87b1-4b1b-99b6-ce90bae932a1",
		Say: "hello",
	}
	fmt.Println(mongo.Upload(&document))
}

func TestDowload(t *testing.T) {
	mongo, err := mongox.Connect("mongodb://admin:admin@127.0.0.1:27017", "parkour")
	if err != nil {
		fmt.Println(err)
	}
	document := Document{
		Id: "2084226a-87b1-4b1b-99b6-ce90bae932a1",
	}
	fmt.Println(mongo.Download(&document), document)
}
