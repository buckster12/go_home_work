package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Front() *listItem
	Back() *listItem
	Len() int
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) Len() int {
	return l.length
}

func (l *list) PushFront(v interface{}) *listItem {
	previousFront := l.front
	l.front = &listItem{Value: v, Prev: previousFront, Next: nil}

	if l.back == nil {
		l.back = l.front
	}

	if previousFront != nil {
		previousFront.Next = l.front
	}

	l.length++
	return l.front
}

func (l *list) PushBack(v interface{}) *listItem {
	previousBack := l.back
	l.back = &listItem{Value: v, Next: previousBack, Prev: nil}

	if l.front == nil {
		l.front = l.back
	}

	if previousBack != nil {
		previousBack.Prev = l.back
	}
	l.length++

	return l.back
}

func (l *list) Remove(i *listItem) {
	switch i {
	case l.back:
		l.back = i.Next
		l.back.Prev = nil
	case l.front:
		l.front = i.Prev
		l.front.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *listItem) {
	if i == l.front {
		return
	}

	prevFront := l.front
	prevFront.Next = i

	l.Remove(i)
	l.length++

	i.Next = nil
	i.Prev = prevFront
	l.front = i
}

type listItem struct {
	Next  *listItem
	Prev  *listItem
	Value interface{}
}

type list struct {
	length int
	front  *listItem
	back   *listItem
}

func NewList() List {
	return &list{}
}
