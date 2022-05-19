/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    lnode
 * @Date:    2022/3/30 16:46
 * @package: lnode
 * @Version: x.x.x
 *
 * @Description: xxx
 *
 */

package lnode

type Head struct {
	First *Node
	Len   int
}

type Node struct {
	Data interface{}
	Next *Node
}

func (l *Node) Link(node *Node) {
	if l.Next == nil {
		l.Next = node
		return
	}
	cur := l.Next
	for cur.Next != nil {
		cur = cur.Next
	}
	cur.Next = node
}

func (l *Node) Push(data interface{}) {
	if l.Next == nil {
		l.Next = &Node{Data: data}
		return
	}
	cur := l.Next
	for cur.Next != nil {
		cur = cur.Next
	}
	cur.Next = &Node{Data: data}
}

func (l *Node) Pop() interface{} {
	if l.Next == nil {
		return l.Data
	}
	cur := l.Next
	var pre *Node
	for cur.Next != nil {
		pre = cur
		cur = cur.Next
	}
	pre.Next = nil
	return cur.Data
}

func (l *Node) Reverse() {
	if l.Next == nil {
		return
	}
	var cur, next *Node
	cur = l.Next.Next
	l.Next.Next = nil
	for cur != nil {
		next = cur.Next
		cur.Next = l.Next
		l.Next = cur
		cur = next
	}
}

func (l *Node) RemoveDup() {
	if l.Next == nil {
		return
	}
	tmp := map[interface{}]struct{}{}
	cur := l
	for cur != nil && cur.Next != nil {
		if _, ok := tmp[cur.Next.Data]; ok {
			cur.Next = cur.Next.Next
		} else {
			tmp[cur.Next.Data] = struct{}{}
			cur = cur.Next
		}
	}
	return
}

func (l *Node) ToArray() []interface{} {
	cur := l.Next
	var reply []interface{}
	for cur != nil {
		reply = append(reply, cur.Data)
		cur = cur.Next
	}
	return reply
}
