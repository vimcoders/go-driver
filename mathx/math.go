package mathx

import "math"

func Sum(values ...int32) (int32, bool) {
	var total int64
	for i := 0; i < len(values); i++ {
		total += int64(values[i])
		if total > math.MaxInt32 {
			return 0, false
		}
	}
	return int32(total), true
}
