package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

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
	l.muFront.Lock()
	defer l.muFront.Unlock()
	return l.front
}

func (l *list) Back() *listItem {
	l.muBack.Lock()
	defer l.muBack.Unlock()
	return l.back
}

func (l *list) Len() int {
	l.muLength.Lock()
	defer l.muLength.Unlock()
	return l.length
}

func (l *list) PushFront(v interface{}) *listItem {
	l.muFront.Lock()
	defer l.muFront.Unlock()
	l.muBack.Lock()
	defer l.muBack.Unlock()
	l.muLength.Lock()
	defer l.muLength.Unlock()

	previousFront := l.front
	l.front = &listItem{Value: v, Prev: previousFront, Next: nil}

	// if we insert the first value (Len() was 0)
	// set front = back, as we have only one value
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
	l.muFront.Lock()
	defer l.muFront.Unlock()
	l.muBack.Lock()
	defer l.muBack.Unlock()
	l.muLength.Lock()
	defer l.muLength.Unlock()

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
	l.muLength.Lock()
	defer l.muLength.Unlock()
	l.muFront.Lock()
	defer l.muFront.Unlock()
	l.muBack.Lock()
	defer l.muBack.Unlock()

	switch i {
	case l.back:
		l.back = i.Next
	case l.front:
		l.front = i.Prev
	default:
	}
	if i.Prev != nil { // not back
		i.Prev.Next = i.Next
	}
	if i.Next != nil { // not front
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *listItem) {
	l.muFront.Lock()
	defer l.muFront.Unlock()
	l.muBack.Lock()
	defer l.muBack.Unlock()

	// already front
	if i == l.front {
		return
	}

	// if back is moving to front
	//if i.Prev == nil {
	if i == l.back {
		l.back = l.back.Next
	} else { // not back
		i.Prev.Next = i.Next
	}
	i.Next.Prev = i.Prev

	prevFront := l.front
	prevFront.Next = i

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
	muLength sync.Mutex
	length   int

	muFront sync.Mutex
	front   *listItem

	muBack sync.Mutex
	back   *listItem
}

func NewList() List {
	return &list{}
}
