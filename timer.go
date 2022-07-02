package driver

import (
	"container/heap"
	"sync"
	"time"
)

// An Item is something we manage in a priority queue.
type Priority struct {
	index    int // The index of the item in the heap.
	Priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	Callback func(priority int)
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Priority

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	p := x.(*Priority)
	p.index = n
	*pq = append(*pq, p)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type WheelTimer struct {
	Id    int64
	Unix  int64 // 时间轮的当前时间
	wheel []heap.Interface
	sync.RWMutex
}

func NewWheelTimer(twNumber int) *WheelTimer {
	if twNumber <= 0 {
		return nil
	}
	tw := &WheelTimer{
		wheel: make([]heap.Interface, twNumber),
		Unix:  time.Now().Unix(),
	}
	for i := 0; i < len(tw.wheel); i++ {
		tw.wheel[i] = &PriorityQueue{}
	}
	return tw
}

func (t *WheelTimer) Push(p *Priority) {
	t.Lock()
	defer t.Unlock()
	h := t.wheel[p.Priority%len(t.wheel)]
	heap.Push(h, p)
}

func (t *WheelTimer) Tick() {
	t.RLock()
	defer t.RUnlock()
	t.Unix++
	h := t.wheel[t.Unix%int64(len(t.wheel))]
	for {
		if h.Len() <= 0 {
			break
		}
		p := heap.Pop(h)
		if priority, ok := p.(*Priority); ok {
			priority.Callback(priority.Priority)
		}
	}
}

func (tw *WheelTimer) Run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		tw.Tick()
	}
}
