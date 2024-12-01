package mongox

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/vimcoders/go-driver/driver"

	"github.com/vimcoders/go-driver/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Queryer[T driver.Document] struct {
	*mongo.Database
}

func Query[T driver.Document](x *Mongo, filter interface{}) (docments []T, err error) {
	return WithQuery[T](x).Query(filter)
}

func WithQuery[T driver.Document](mongo *Mongo) *Queryer[T] {
	return &Queryer[T]{Database: mongo.Database}
}

func (x Queryer[T]) Query(filter interface{}) (docments []T, err error) {
	var doc T
	c := x.Collection(doc.DocumentName())
	if c == nil {
		return nil, errors.New("no collection")
	}
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(context.Background(), &docments); err != nil {
		return nil, err
	}
	return docments, nil
}

type Mongo struct {
	*mongo.Database
}

func (x *Mongo) Close() error {
	return x.Client().Disconnect(context.Background())
}

func (x *Mongo) Query(document string, filter interface{}, results interface{}, opts ...*options.FindOptions) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	c := x.Collection(document)
	if c == nil {
		return errors.New("no collection")
	}
	cur, err := c.Find(context.Background(), filter, opts...)
	if err != nil {
		return err
	}
	return cur.All(context.Background(), results)
}

func (x *Mongo) Insert(documents ...driver.Document) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(documents) <= 0 {
		return nil
	}
	for i := 0; i < len(documents); i++ {
		c := x.Collection(documents[i].DocumentName())
		if c == nil {
			continue
		}
		if _, err := c.InsertOne(context.Background(), documents[i]); err != nil {
			return err
		}
	}
	return nil
}

func (x *Mongo) Delete(documents ...driver.Document) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(documents) <= 0 {
		return nil
	}
	for i := 0; i < len(documents); i++ {
		c := x.Collection(documents[i].DocumentName())
		if c == nil {
			continue
		}
		if _, err := c.DeleteOne(context.Background(), documents[i]); err != nil {
			return err
		}
	}
	return nil
}

func (x *Mongo) Update(documents ...driver.Document) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(documents) <= 0 {
		return nil
	}
	for i := 0; i < len(documents); i++ {
		c := x.Collection(documents[i].DocumentName())
		if c == nil {
			continue
		}
		if _, err := c.UpdateByID(context.Background(), documents[i].DocumentId(), bson.M{"$set": documents[i]}); err != nil {
			return err
		}
	}
	return nil
}

func (x *Mongo) Upsert(filter interface{}, document driver.Document) (err error) {
	ctx := context.Background()
	opt := options.Update().SetUpsert(true)
	c := x.Collection(document.DocumentName())
	r, err := c.UpdateOne(ctx, filter, bson.M{"$set": document}, opt)
	log.Debugf("Upsert collection:%s, id:%v, result:%+v, error:%v",
		document.DocumentName(), document.DocumentId(), r, err)
	return err
}

func Connect(host, db string) (*Mongo, error) {
	opts := options.Client().ApplyURI(host)
	opts.SetMaxPoolSize(8)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	mongodb := client.Database(db)
	if mongodb == nil {
		panic("no db")
	}
	return &Mongo{
		Database: mongodb,
	}, nil
}

func (x *Mongo) UpdateByID(documents ...driver.Document) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(documents) <= 0 {
		return nil
	}
	for i := 0; i < len(documents); i++ {
		c := x.Collection(documents[i].DocumentName())
		if c == nil {
			continue
		}
		if _, err := c.UpdateByID(context.Background(), documents[i].DocumentId(), documents[i]); err != nil {
			return err
		}
	}
	return nil
}

func (x *Mongo) Upload(doc driver.Document) error {
	b, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	bucket, err := gridfs.NewBucket(x.Database, options.GridFSBucket().SetName(doc.DocumentName()))
	if err != nil {
		return err
	}
	if err := bucket.Delete(doc.DocumentId()); err != nil {
		log.Error(err.Error())
	}
	return bucket.UploadFromStreamWithID(doc.DocumentId(), doc.DocumentId()+".report", bytes.NewBuffer(b))
}

func (x *Mongo) Download(doc driver.Document) error {
	bucket, err := gridfs.NewBucket(x.Database, options.GridFSBucket().SetName(doc.DocumentName()))
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(nil)
	if _, err := bucket.DownloadToStream(doc.DocumentId(), buffer); err != nil {
		return err
	}
	if err := json.Unmarshal(buffer.Bytes(), doc); err != nil {
		return err
	}
	return nil
}
