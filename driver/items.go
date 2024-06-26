// 不允许调用标准库外的包，防止循环引用
package driver

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Item struct {
	Id    int32 `json:"id"`
	Count int32 `json:"count"`
}

// func (x *Item) ToMessage() *pb.Item {
// 	return &pb.Item{
// 		Type:  x.Id,
// 		Count: x.Count,
// 	}
// }

type Items []*Item

func (x *Items) Add(items ...*Item) {
	for i := 0; i < len(items); i++ {
		if items[i].Count <= 0 {
			continue
		}
		x.add(items[i])
	}
}

func (x *Items) add(item *Item) {
	for _, v := range *x {
		if v.Id != item.Id {
			continue
		}
		//v.Count = mathx.Sum(v.Count, item.Count)
		return
	}
	*x = append(*x, &Item{
		Id:    item.Id,
		Count: item.Count,
	})
}

func (x Items) Sub(items ...*Item) {
	var values Items
	values.Add(items...)
	for i := 0; i < len(values); i++ {
		x.sub(values[i])
	}
}

func (x Items) sub(item *Item) {
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

func (x Items) Get(t int32) int32 {
	for i := 0; i < len(x); i++ {
		if x[i].Id != t {
			continue
		}
		return x[i].Count
	}
	return 0
}

// func (x *ItemList) SubWithSafe(items ...*Item) (unenough int32, syncItems []*pb.Item, ok bool) {
// 	var values ItemList
// 	values.Add(items...)
// 	for i := 0; i < len(values); i++ {
// 		if x.Get(values[i].Id) < values[i].Count {
// 			return values[i].Id, nil, false
// 		}
// 	}
// 	for i := 0; i < len(values); i++ {
// 		sync, ok := x.subWithSafe(values[i])
// 		if !ok {
// 			return 0, nil, false
// 		}
// 		syncItems = append(syncItems, &pb.Item{
// 			Type:  values[i].Id,
// 			Count: sync,
// 		})
// 	}
// 	return 0, syncItems, true
// }

// func (x ItemList) subWithSafe(item *Item) (int32, bool) {
// 	for i := 0; i < len(x); i++ {
// 		if x[i].Id != item.Id {
// 			continue
// 		}
// 		if x[i].Count < item.Count {
// 			return 0, false
// 		}
// 		x[i].Count -= item.Count
// 		return x[i].Count, true
// 	}
// 	return 0, false
// }

// func (x ItemList) Multiply(n int32) (clone ItemList) {
// 	if n <= 0 {
// 		return
// 	}
// 	for i := 0; i < len(x); i++ {
// 		count := int64(x[i].Count) * int64(n)
// 		clone = append(clone, &Item{
// 			Id:    x[i].Id,
// 			Count: int32(mathx.Min(count, math.MaxInt32)),
// 		})
// 	}
// 	return clone
// }

// func (x ItemList) Percent(percent int32) (clone ItemList) {
// 	if percent <= 0 || percent > 10000 {
// 		return nil
// 	}
// 	for i := 0; i < len(x); i++ {
// 		clone = append(clone, &Item{
// 			Id:    x[i].Id,
// 			Count: int32(math.Floor(float64(percent) / 10000 * float64(x[i].Count))),
// 		})
// 	}
// 	return clone
// }

// func (x ItemList) ToMessage() (result []*pb.Item) {
// 	for i := 0; i < len(x); i++ {
// 		result = append(result, x[i].ToMessage())
// 	}
// 	return result
// }

// func (x ItemList) Total() (total int64) {
// 	for i := 0; i < len(x); i++ {
// 		total = int64(mathx.Sum(int32(total), x[i].Count))
// 	}
// 	return total
// }

func (x Items) Clone() (copy Items) {
	for i := 0; i < len(x); i++ {
		copy = append(copy, &Item{
			Id:    x[i].Id,
			Count: x[i].Count,
		})
	}
	return copy
}

func (x Items) ToString() (str string) {
	for i := 0; i < len(x); i++ {
		str += fmt.Sprintf("%v,%v;", x[i].Id, x[i].Count)
	}
	return strings.TrimRight(str, ";")
}

func (x *Items) Unmarshal(b string) error {
	for _, v := range strings.Split(b, ";") {
		itemStr := strings.Split(v, ",")
		if len(itemStr) < 2 {
			continue
		}
		t, err := strconv.Atoi(itemStr[0])
		if err != nil {
			return err
		}
		count, err := strconv.Atoi(itemStr[1])
		if err != nil {
			return err
		}
		*x = append(*x, &Item{
			Id:    int32(t),
			Count: int32(count),
		})
	}
	return nil
}
