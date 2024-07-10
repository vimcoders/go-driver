// 不允许调用标准库外的包，防止循环引用
package driver

// type Distinct[T comparable] []T

// func (x Distinct[T]) Get(e T) (int, bool) {
// 	for i := 0; i < len(x); i++ {
// 		if x[i] == e {
// 			return i + 1, true
// 		}
// 	}
// 	return 0, false
// }

// func (x Distinct[T]) Clone() (distinct Distinct[T]) {
// 	for i := 0; i < len(x); i++ {
// 		distinct = append(distinct, x[i])
// 	}
// 	return distinct
// }

// func (x *Distinct[T]) Push(values ...T) bool {
// 	if len(values) <= 0 {
// 		return false
// 	}
// 	distinct, ok := NewDistinct[T](values...)
// 	if !ok {
// 		return false
// 	}
// 	old := *x
// 	for i := 0; i < len(old); i++ {
// 		if _, ok := distinct.Get(old[i]); ok {
// 			return false
// 		}
// 		distinct = append(distinct, old[i])
// 	}
// 	*x = distinct
// 	return true
// }

// func (x *Distinct[T]) Delete(e T) bool {
// 	old := *x
// 	for i := 0; i < len(old); i++ {
// 		if old[i] != e {
// 			continue
// 		}
// 		*x = append(old[:i], old[i+1:]...)
// 		return true
// 	}
// 	return false
// }

// func (x Distinct[T]) Compare(set Distinct[T]) bool {
// 	if len(x) != len(set) {
// 		return false
// 	}
// 	for i := 0; i < len(x); i++ {
// 		if x[i] != set[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func NewDistinct[T comparable](values ...T) (Distinct[T], bool) {
// 	var distinct Distinct[T]
// 	for i := 0; i < len(values); i++ {
// 		if _, ok := distinct.Get(values[i]); ok {
// 			return nil, false
// 		}
// 		distinct = append(distinct, values[i])
// 	}
// 	return distinct, true
// }
