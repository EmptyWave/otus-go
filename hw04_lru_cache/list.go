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
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.len > 0 {
		item.Next = l.front
		l.front.Prev = item
	} else {
		l.back = item
	}

	l.front = item
	l.len++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.len > 0 {
		item.Prev = l.back
		l.back.Next = item
	} else {
		l.front = item
	}

	l.back = item
	l.len++

	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)

	i.Prev = nil
	i.Next = l.front

	l.front.Prev = i
	l.front = i
	l.len++
}
