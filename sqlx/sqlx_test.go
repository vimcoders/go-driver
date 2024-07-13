package sqlx_test

import (
	"fmt"
	"go-driver/sqlx"
	"math/rand"
	"testing"
	"time"
)

// go test -v -bench=BenchmarkQuery
func BenchmarkQuery(b *testing.B) {
	type Account struct {
		UserId   uint64 `gorm:"primarykey"`
		Passport string `gorm:"unique"`
		Pwd      string
		Created  time.Time `gorm:"comment:创建时间"`
	}
	client, err := sqlx.Dial("root:root@tcp(127.0.0.1:3306)/proxy?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		query := client.Where("user_id", rand.Uint64())
		var account Account
		if err := query.Query(&account); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -v -bench=BenchmarkInsert
func BenchmarkInsert(b *testing.B) {
	type Account struct {
		UserId   uint64 `gorm:"primarykey"`
		Passport string `gorm:"unique"`
		Pwd      string
		Created  time.Time `gorm:"comment:创建时间"`
	}
	client, err := sqlx.Dial("root:root@tcp(127.0.0.1:3306)/proxy?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := &Account{
			Passport: fmt.Sprintf("%d", rand.Int63()),
			Pwd:      fmt.Sprintf("%d", rand.Int63()),
			Created:  time.Now(),
		}
		if err := client.Insert(&account); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -v -bench=BenchmarkDelete
func BenchmarkDelete(b *testing.B) {
	type Account struct {
		UserId   uint64 `gorm:"primarykey"`
		Passport string `gorm:"unique"`
		Pwd      string
		Created  time.Time `gorm:"comment:创建时间"`
	}
	client, err := sqlx.Dial("root:root@tcp(127.0.0.1:3306)/proxy?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		delete := client.Where("user_id", rand.Uint64())
		if err := delete.Delete(&Account{}); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -v -bench=BenchmarkUpdate
func BenchmarkUpdate(b *testing.B) {
	type Account struct {
		UserId   uint64 `gorm:"primarykey"`
		Passport string `gorm:"unique"`
		Pwd      string
		Created  time.Time `gorm:"comment:创建时间"`
	}
	client, err := sqlx.Dial("root:root@tcp(127.0.0.1:3306)/proxy?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := &Account{
			Passport: fmt.Sprintf("%d", rand.Int63()),
			Pwd:      fmt.Sprintf("%d", rand.Int63()),
			Created:  time.Now(),
		}
		update := client.Where("user_id", i)
		if err := update.Update(account); err != nil {
			b.Error(err)
			return
		}
	}
}
