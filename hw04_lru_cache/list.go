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

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) pushFrontItem(i *ListItem) *ListItem {
	if l.len == 0 {
		l.front = i
		l.back = l.front
	} else {
		l.front.Prev = i
		l.front.Prev.Next = l.front
		l.front = l.front.Prev
	}
	l.len++
	return l.front
}

func (l *list) pushBackItem(i *ListItem) *ListItem {
	if l.len == 0 {
		l.front = i
		l.back = l.front
	} else {
		l.back.Next = i
		l.back.Next.Prev = l.back
		l.back = l.back.Next
	}
	l.len++
	return l.back
}

func (l *list) pop(i *ListItem) *ListItem {
	remItem := l.front
	for remItem != nil && remItem.Value != i.Value {
		remItem = remItem.Next // not in circle for linter
	}
	if remItem.Next != nil {
		remItem.Next.Prev = remItem.Prev
	}
	if remItem.Prev != nil {
		remItem.Prev.Next = remItem.Next
	}
	if remItem == l.front {
		l.front = remItem.Next
	}
	if remItem == l.back {
		l.back = remItem.Prev
	}
	l.len--
	return remItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	return l.pushFrontItem(newItem)
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	return l.pushBackItem(newItem)
}

func (l *list) Remove(i *ListItem) {
	l.pop(i)
}

func (l *list) MoveToFront(i *ListItem) {
	l.pushFrontItem(l.pop(i))
}

func NewList() List {
	return new(list)
}
