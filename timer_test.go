package driver

import (
	"fmt"
	"testing"
	"time"
)

func TestWheelTimer(t *testing.T) {
	wheel := NewWheelTimer(10)
	go func() {
		for {
			wheel.Push(&Priority{
				Priority: int(time.Now().Add(time.Second * 5).Unix()),
				Callback: func(priority int) {
					fmt.Println(time.Unix(int64(priority), 0))
				},
			})
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		for {
			wheel.Push(&Priority{
				Priority: int(time.Now().Add(time.Second * 5).Unix()),
				Callback: func(priority int) {
					fmt.Println(time.Unix(int64(priority), 0))
				},
			})
			time.Sleep(time.Millisecond * 100)
		}
	}()
	wheel.Run()
}
