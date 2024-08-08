// 不允许调用标准库外的包，防止循环引用
package driver

type Task struct {
	Id    int32
	Count int32
}

type TaskList []*Task

func (x TaskList) Get(taskId int32) (*Task, bool) {
	for i := 0; i < len(x); i++ {
		if x[i].Id != taskId {
			continue
		}
		return x[i], true
	}
	return nil, false
}
