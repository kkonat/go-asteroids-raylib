package main

type ListEl[T any] struct {
	Prev, Next *ListEl[T]
	Value      T
}

type List[T any] struct {
	Head, Tail *ListEl[T]
	Len        int
}

func (dll *List[T]) AppendEl(el *ListEl[T]) {

	if dll.Head == nil && dll.Tail == nil { // empty list
		dll.Head = el
		dll.Tail = el
	} else { // head -> [data0] ... [datan] <- tail
		dll.Tail.Next = el //update dll
		el.Prev = dll.Tail // update tll
		dll.Tail = el      // update tail
	}
	dll.Len++
}
func (dll *List[T]) AppendVal(val T) *ListEl[T] {
	ne := &ListEl[T]{Value: val}

	if dll.Head == nil && dll.Tail == nil { // empty list
		dll.Head = ne
		dll.Tail = ne
	} else { // head -> [data0] ... [datan] <- tail
		dll.Tail.Next = ne //update dll
		ne.Prev = dll.Tail // update tll
		dll.Tail = ne      // update tail
	}
	dll.Len++
	return ne
}

func (dll *List[T]) Delete(el *ListEl[T]) bool {

	if (dll.Head == nil && dll.Tail == nil) || (el == nil) {
		return false
	}
	if (el.Next == nil) && (el.Prev == nil) { //single element
		dll.Head, dll.Tail = nil, nil
	} else if el.Prev == nil { // at head
		dll.Head = el.Next
		el.Next.Prev = nil
	} else if el.Next == nil { //at tail
		el.Prev.Next = nil
		dll.Tail = el.Prev
	} else {
		el.Prev.Next = el.Next
		el.Next.Prev = el.Prev
	}
	dll.Len--
	el = nil
	return true
}
func (dll *List[T]) Clear() {
	for dll.Len > 0 {
		dll.Delete(dll.Tail)
	}
}
func (dll *List[T]) Iter() func() (*ListEl[T], bool) {
	ptr := dll.Head
	return func() (*ListEl[T], bool) {
		if ptr == dll.Tail {
			return nil, false
		}
		retVal := ptr
		ptr = ptr.Next
		return retVal, true
	}
}
