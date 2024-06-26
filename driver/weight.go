// 不允许调用标准库外的包，防止循环引用
package driver

// type Weight struct {
// 	Id     int
// 	Weight int
// }

// type WeightList []*Weight

// func (x WeightList) Rand() (int, bool) {
// 	var total int
// 	for _, v := range x {
// 		if v.Id <= 0 {
// 			return 0, false
// 		}
// 		if v.Weight <= 0 {
// 			return 0, false
// 		}
// 		total += v.Weight
// 	}
// 	if total <= 0 {
// 		return 0, false
// 	}
// 	rand.New(rand.NewSource(time.Now().UnixMicro()))
// 	randomNum := rand.Intn(total) + 1
// 	sort.Slice(x, func(i, j int) bool {
// 		return x[i].Weight < x[j].Weight
// 	})
// 	for _, v := range x {
// 		randomNum -= v.Weight
// 		if randomNum > 0 {
// 			continue
// 		}
// 		return v.Id, true
// 	}
// 	return 0, false
// }

// func (x *WeightList) Unmarshal(b string) error {
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
// 		*x = append(*x, &Weight{
// 			Id:     t,
// 			Weight: count,
// 		})
// 	}
// 	return nil
// }
