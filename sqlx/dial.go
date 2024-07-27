package sqlx

import (
	"go-driver/log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Dial(host string) (Client, error) {
	log.Debugf("host := %s", host)
	db, err := gorm.Open(mysql.Open(host), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	//db, err := gorm.Open(mysql.Open(host), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(16)
	sqlDB.SetConnMaxLifetime(59 * time.Second)
	return &XClient{DB: db}, nil
}
