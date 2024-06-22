// https://github.com/HDT3213/godis.git
package sortedset

import (
	"math/rand"
	"time"
)

const MaxLevel = 32
const p = 0.5

type Node struct {
	value  uint32
	levels []*Level // 索引节点,index=0是基础链表
}

type Level struct {
	next *Node
}

type SkipList struct {
	header *Node  // 表头节点
	length uint32 // 原始链表的长度，表头节点不计入
	height uint32 // 最高的节点的层数
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: NewNode(MaxLevel, 0),
		length: 0,
		height: 1,
	}
}

func NewNode(level, value uint32) *Node {
	node := new(Node)
	node.value = value
	node.levels = make([]*Level, level)

	for i := 0; i < len(node.levels); i++ {
		node.levels[i] = new(Level)
	}
	return node
}

func (sl *SkipList) randomLevel() int {
	level := 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for r.Float64() < p && level < MaxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Add(value uint32) bool {
	if value <= 0 {
		return false
	}
	update := make([]*Node, MaxLevel)
	// 每一次循环都是一次寻找有序单链表的插入过程
	tmp := sl.header
	for i := int(sl.height) - 1; i >= 0; i-- {
		// 每次循环不重置 tmp，直接从上一层确认的节点开始向下一层查找
		for tmp.levels[i].next != nil && tmp.levels[i].next.value < value {
			tmp = tmp.levels[i].next
		}

		// 避免插入相同的元素
		if tmp.levels[i].next != nil && tmp.levels[i].next.value == value {
			return false
		}

		update[i] = tmp
	}

	level := sl.randomLevel()
	node := NewNode(uint32(level), value)
	// fmt.Printf("level:%v,value:%v\n", level, value)

	if uint32(level) > sl.height {
		sl.height = uint32(level)
	}

	for i := 0; i < level; i++ {

		// 说明新节点层数超过了跳表当前的最高层数，此时将头节点对应层数的后继节点设置为新节点
		if update[i] == nil {
			sl.header.levels[i].next = node
			continue
		}
		// 普通的插入链表操作，修改新节点的后继节点为前一个节点的后继节点，修改前一个节点的后继节点为新节点
		node.levels[i].next = update[i].levels[i].next
		update[i].levels[i].next = node
	}

	sl.length++
	return true
}

func (sl *SkipList) Delete(value uint32) bool {
	var node *Node
	last := make([]*Node, sl.height)
	tmp := sl.header
	for i := int(sl.height) - 1; i >= 0; i-- {

		for tmp.levels[i].next != nil && tmp.levels[i].next.value < value {
			tmp = tmp.levels[i].next
		}

		last[i] = tmp
		// 拿到 value 对应的 node
		if tmp.levels[i].next != nil && tmp.levels[i].next.value == value {
			node = tmp.levels[i].next
		}
	}

	// 没有找到 value 对应的 node
	if node == nil {
		return false
	}

	// 找到所有前置节点后需要删除node
	for i := 0; i < len(node.levels); i++ {
		last[i].levels[i].next = node.levels[i].next
		node.levels[i].next = nil
	}

	// 重定向跳表高度
	for i := 0; i < len(sl.header.levels); i++ {
		if sl.header.levels[i].next == nil {
			sl.height = uint32(i)
			break
		}
	}

	sl.length--

	return true
}
