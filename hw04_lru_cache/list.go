package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	first *ListItem
	last  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.len == 0 {
		l.first = &ListItem{v, nil, nil}
		l.last = l.first
	} else {
		newItem := &ListItem{v, l.first, nil}
		l.first.Prev = newItem
		l.first = newItem
	}
	l.len++
	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.len == 0 {
		l.first = &ListItem{v, nil, nil}
		l.last = l.first
	} else {
		newItem := &ListItem{v, nil, l.last}
		l.last.Next = newItem
		l.last = newItem
	}
	l.len++
	return l.last
}

func (l *list) Remove(i *ListItem) {
	cur := l.first
	for cur != nil {
		if cur == i {
			prev := cur.Prev
			if cur == l.last {
				l.last = prev
			}
			next := cur.Next
			if cur == l.first {
				l.first = next
			}
			if prev != nil {
				prev.Next = next
			}
			if next != nil {
				next.Prev = prev
			}
			l.len--
			break
		}
		cur = cur.Next
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
