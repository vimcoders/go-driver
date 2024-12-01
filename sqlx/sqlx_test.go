package sqlx_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/vimcoders/go-driver/sqlx"
)

// go test -v -bench=BenchmarkRegister
func BenchmarkRegister(b *testing.B) {
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
		if err := client.Register(&Account{}); err != nil {
			b.Error(err)
			return
		}
	}
}

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
	if err := client.Register(&Account{}); err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := &Account{
			UserId: uint64(rand.Uint32()),
		}
		if err := client.Query(&account); err != nil {
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
	if err := client.Register(&Account{}); err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := &Account{
			UserId:   uint64(rand.Uint32()),
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
	if err := client.Register(&Account{}); err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := Account{
			UserId:   uint64(rand.Uint32()),
			Passport: fmt.Sprintf("%d", rand.Int63()),
			Pwd:      fmt.Sprintf("%d", rand.Int63()),
			Created:  time.Now(),
		}
		if err := client.Delete(&account); err != nil {
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
	if err := client.Register(&Account{}); err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := &Account{
			UserId:   uint64(rand.Uint32()),
			Passport: fmt.Sprintf("%d", rand.Int63()),
			Pwd:      fmt.Sprintf("%d", rand.Int63()),
			Created:  time.Now(),
		}
		if err := client.Update(account); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -v -bench=BenchmarkReplace
func BenchmarkReplace(b *testing.B) {
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
	if err := client.Register(&Account{}); err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		account := &Account{
			UserId:   uint64(i + 1),
			Passport: fmt.Sprintf("%d", rand.Int63()),
			Pwd:      fmt.Sprintf("%d", rand.Int63()),
			Created:  time.Now(),
		}
		if err := client.Replace(&account); err != nil {
			b.Error(err)
			return
		}
	}
}
