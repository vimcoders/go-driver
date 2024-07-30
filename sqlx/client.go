package sqlx

import (
	"errors"

	"gorm.io/gorm"
)

type Client interface {
	Register(...interface{}) error
	Insert(...interface{}) error
	Delete(...interface{}) error
	Query(interface{}) error
	Update(...interface{}) error
	Replace(values ...interface{}) (err error)
	Where(string, ...interface{}) Client
	Close() error
}

type XClient struct {
	*gorm.DB
}

func (x *XClient) Where(query string, args ...interface{}) Client {
	return &XClient{DB: x.DB.Where(query, args...)}
}

func (x *XClient) Register(values ...interface{}) error {
	return x.Migrator().AutoMigrate(values...)
}

func (x *XClient) Insert(values ...interface{}) (err error) {
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	tx := x.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for i := 0; i < len(values); i++ {
		if t := tx.Create(values[i]); t != nil && t.Error != nil {
			return t.Error
		}
	}
	if t := tx.Commit(); t != nil && t.Error != nil {
		return t.Error
	}
	return nil
}

func (x *XClient) Delete(values ...interface{}) (err error) {
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	tx := x.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for i := 0; i < len(values); i++ {
		if t := tx.Delete(values[i]); t != nil && t.Error != nil {
			return t.Error
		}
	}
	if t := tx.Commit(); t != nil && t.Error != nil {
		return t.Error
	}
	return nil
}

func (x *XClient) Query(dest interface{}) error {
	if tx := x.Find(dest); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (x *XClient) Update(values ...interface{}) (err error) {
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	tx := x.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for i := 0; i < len(values); i++ {
		if t := tx.Updates(values[i]); t != nil && t.Error != nil {
			return t.Error
		}
	}
	if t := tx.Commit(); t != nil && t.Error != nil {
		return t.Error
	}
	return nil
}

func (x *XClient) Replace(values ...interface{}) (err error) {
	if len(values) <= 0 {
		return errors.New("len(values) <= 0")
	}
	tx := x.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for i := 0; i < len(values); i++ {
		if err := tx.Save(values[i]).Error; err != nil {
			return err
		}
	}
	if t := tx.Commit(); t != nil && t.Error != nil {
		return t.Error
	}
	return nil
}

func (x *XClient) Close() error {
	db, err := x.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
