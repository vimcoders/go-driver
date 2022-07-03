package driver

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestWheelTimer(t *testing.T) {
	wheel := NewWheelTimer(1000)
	go func() {
		for {
			wheel.Push(&Priority{
				Priority: int(time.Now().Unix()) + rand.Intn(10),
				Callback: func(priority int) {
					fmt.Println(time.Unix(int64(priority), 0), 1)
				},
			})
			time.Sleep(time.Millisecond)
		}
	}()
	go func() {
		for {
			wheel.Push(&Priority{
				Priority: int(time.Now().Unix()) + rand.Intn(10),
				Callback: func(priority int) {
					fmt.Println(time.Unix(int64(priority), 0), 2)
				},
			})
			time.Sleep(time.Millisecond)
		}
	}()
	wheel.Run()
}
