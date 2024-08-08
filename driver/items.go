// 不允许调用标准库外的包，防止循环引用
package driver

import (
	"go-driver/mathx"
	"math"
)

type Item struct {
	Id    int32 `json:"id"`
	Count int32 `json:"count"`
}

type ItemList []*Item

func (x *ItemList) Add(items ...*Item) bool {
	var values ItemList
	for i := 0; i < len(items); i++ {
		if ok := values.add(items[i]); !ok {
			return false
		}
	}
	for i := 0; i < len(values); i++ {
		if values[i].Count <= 0 {
			return false
		}
	}
	for i := 0; i < len(values); i++ {
		if ok := x.add(values[i]); !ok {
			return false
		}
	}
	return true
}

func (x *ItemList) add(item *Item) bool {
	for _, v := range *x {
		if v.Id != item.Id {
			continue
		}
		if count, ok := mathx.Sum(v.Count, item.Count); ok {
			v.Count = count
			return true
		}
		return false
	}
	*x = append(*x, &Item{
		Id:    item.Id,
		Count: item.Count,
	})
	return true
}

func (x *ItemList) Sub(items ...*Item) {
	var values ItemList
	values.Add(items...)
	for i := 0; i < len(values); i++ {
		x.sub(values[i])
	}
}

func (x ItemList) sub(item *Item) {
	for i := 0; i < len(x); i++ {
		if x[i].Id != item.Id {
			continue
		}
		if int64(x[i].Count)-int64(item.Count) < math.MinInt32 {
			x[i].Count = math.MinInt32
			continue
		}
		x[i].Count -= item.Count
	}
}

func (x ItemList) Get(t int32) (*Item, bool) {
	for i := 0; i < len(x); i++ {
		if x[i].Id != t {
			continue
		}
		return x[i], true
	}
	return nil, false
}

func (x ItemList) Clone() (copy ItemList) {
	for i := 0; i < len(x); i++ {
		copy = append(copy, &Item{
			Id:    x[i].Id,
			Count: x[i].Count,
		})
	}
	return copy
}
