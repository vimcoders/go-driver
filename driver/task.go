// 不允许调用标准库外的包，防止循环引用
package driver

// import "parkour/pb"

type Task struct {
	Id    int32
	Count int32
}

// func (x *Task) ToMessage() *pb.Task {
// 	return &pb.Task{
// 		Id:    x.Id,
// 		Count: x.Count,
// 	}
// }

type Tasks []*Task

// func (x TaskList) ToMessage() (result []*pb.Task) {
// 	for i := 0; i < len(x); i++ {
// 		result = append(result, x[i].ToMessage())
// 	}
// 	return result
// }

// func (x TaskList) ToPbTaskList() (result []*pb.Task) {
// 	for i := 0; i < len(x); i++ {
// 		result = append(result, x[i].ToMessage())
// 	}
// 	return result
// }

func (x Tasks) Get(taskId int32) (int32, bool) {
	for i := 0; i < len(x); i++ {
		if x[i].Id != taskId {
			continue
		}
		return x[i].Count, true
	}
	return 0, false
}
