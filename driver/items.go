package driver

// import (
// 	"fmt"
// 	"math"
// 	"parkour/lib/mathx"
// 	"parkour/pb"
// 	"strconv"
// 	"strings"
// )

// const (
// 	ITEM_EXP                   = -8
// 	ITEM_GOLD                  = -7    // 金币
// 	ITEM_DIAMOND               = -6    // 钻石
// 	ITEM_CROWN                 = -3    // 皇冠
// 	ITEM_TROPHIES              = -4    // 奖杯
// 	ITEM_SPORTS_TICKET         = -2    // 竞技场门票
// 	ITEM_SPORTS_MASTER_TICKET  = -15   // 竞技场大师门票
// 	ITEM_SPORTS_ADVANCE_TICKET = -14   // 竞技场高级门票
// 	ITEM_SPORTS_MIDDLE_TICKET  = -13   // 竞技场中级门票
// 	ITEM_CASUAL_TICKET         = -10   // 休闲场门票
// 	ITEM_SPORTS_WEEKLY_POINT   = -12   // 周竞技积分
// 	ITEM_SPORTS_MONTHLY_POINT  = -11   // 月竞技积分
// 	ITEM_DAILY_POINT           = 31000 // 日常活跃点
// 	ITEM_WEEK_POINT            = 31001 // 周常活跃点
// 	ITEM_ACHIEVEMENT_POINT     = 31002 // 成就活跃点
// 	ITEM_COLLECT_THIRD_NINE    = 60003 // 集九：第3个九
// 	ITEM_COLLECT_FIFTH_NINE    = 60005 // 集九：第5个九
// )

// type Item struct {
// 	Id    int32 `json:"id"`
// 	Count int32 `json:"count"`
// }

// func (x *Item) ToMessage() *pb.Item {
// 	return &pb.Item{
// 		Type:  x.Id,
// 		Count: x.Count,
// 	}
// }

// type ItemList []*Item

// func (x *ItemList) Add(items ...*Item) ItemList {
// 	for _, v := range items {
// 		if v.Count <= 0 {
// 			continue
// 		}
// 		x.add(v)
// 	}
// 	return *x
// }

// func (x ItemList) Sub(items ...*Item) {
// 	var values ItemList
// 	values.Add(items...)
// 	for i := 0; i < len(values); i++ {
// 		x.sub(values[i])
// 	}
// }

// func (x ItemList) sub(item *Item) {
// 	for i := 0; i < len(x); i++ {
// 		if x[i].Id != item.Id {
// 			continue
// 		}
// 		if int64(x[i].Count)-int64(item.Count) < math.MinInt32 {
// 			x[i].Count = math.MinInt32
// 			continue
// 		}
// 		x[i].Count -= item.Count
// 	}
// }

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

// func (x *ItemList) add(item *Item) {
// 	for _, v := range *x {
// 		if v.Id != item.Id {
// 			continue
// 		}
// 		v.Count = mathx.Sum(v.Count, item.Count)
// 		return
// 	}
// 	*x = append(*x, &Item{
// 		Id:    item.Id,
// 		Count: item.Count,
// 	})
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

// func (x ItemList) Get(t int32) int32 {
// 	for _, v := range x {
// 		if v.Id != t {
// 			continue
// 		}
// 		return v.Count
// 	}
// 	return 0
// }

// func (x ItemList) Total() (total int64) {
// 	for i := 0; i < len(x); i++ {
// 		total = int64(mathx.Sum(int32(total), x[i].Count))
// 	}
// 	return total
// }

// func (x ItemList) Clone() (copy ItemList) {
// 	for i := 0; i < len(x); i++ {
// 		copy = append(copy, &Item{
// 			Id:    x[i].Id,
// 			Count: x[i].Count,
// 		})
// 	}
// 	return copy
// }

// func (x ItemList) ToString() (str string) {
// 	for i := 0; i < len(x); i++ {
// 		str += fmt.Sprintf("%v,%v;", x[i].Id, x[i].Count)
// 	}
// 	return strings.TrimRight(str, ";")
// }

// func (x *ItemList) Unmarshal(b string) error {
// 	for _, v := range strings.Split(b, ";") {
// 		itemStr := strings.Split(v, ",")
// 		if len(itemStr) < 2 {
// 			continue
// 		}
// 		t, err := strconv.Atoi(itemStr[0])
// 		if err != nil {
// 			return err
// 		}
// 		count, err := strconv.Atoi(itemStr[1])
// 		if err != nil {
// 			return err
// 		}
// 		*x = append(*x, &Item{
// 			Id:    int32(t),
// 			Count: int32(count),
// 		})
// 	}
// 	return nil
// }
