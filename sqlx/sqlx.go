package sqlx

import (
	"time"

	"go-driver/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	*gorm.DB
}

func (x *Client) Register(values ...interface{}) {
	x.Migrator().AutoMigrate(values...)
}

func (x *Client) Insert(values ...interface{}) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(values) <= 0 {
		return
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

func (x *Client) Delete(values ...interface{}) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(values) <= 0 {
		return
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

func (x *Client) Update(values ...interface{}) (err error) {
	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()
	if len(values) <= 0 {
		return
	}
	tx := x.Begin()
	defer func() {
		if err != nil {
			log.Error(err.Error())
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

func (x *Client) Close() error {
	db, err := x.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func Connect(host string) (*Client, error) {
	log.Debugf("host:= %s", host)
	db, err := gorm.Open(mysql.Open(host), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	//db, err := gorm.Open(mysql.Open(host), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")
	log.Debugf("数据库地址: %s", host)
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(16)
	sqlDB.SetConnMaxLifetime(59 * time.Second)
	return &Client{DB: db}, nil
}

// func Where[T any](query string, args ...interface{}) (T, error) {
// 	var dest T
// 	dbMysql := Mysql()
// 	if dbMysql == nil {
// 		return dest, errors.New("dbMysql == nil")
// 	}
// 	if tx := dbMysql.Where(query, args...).Find(&dest); tx.Error != nil {
// 		return dest, tx.Error
// 	}
// 	return dest, nil
// }

// func Count[T any](query string, args ...interface{}) (int64, error) {
// 	dbMysql := Mysql()
// 	if dbMysql == nil {
// 		return 0, errors.New("dbMysql == nil")
// 	}
// 	var dest T
// 	var count int64
// 	if tx := dbMysql.Model(dest).Where(query, args...).Count(&count); tx.Error != nil {
// 		return 0, tx.Error
// 	}
// 	return count, nil
// }
