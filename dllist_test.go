package main

import (
	"testing"
)

func TestDll(t *testing.T) {

	dl := List[int]{}
	if dl.Head != dl.Tail {
		t.Error("dl.Head != dl.Tail")
	}
	if dl.Head != nil || dl.Tail != nil {
		t.Error("dl.Head != nil || dl.Tail != nil")
	}

	dl.AppendVal(1)
	if dl.Head != dl.Tail {
		t.Error("dl.Head != dl.Tail")
	}
	if dl.Head.Prev != nil || dl.Head.Next != nil {
		t.Error("dl.Head.Prev != nil || dl.Head.Next !=nil")
	}
	el := dl.Head
	dl.Delete(el)
	if dl.Head != nil || dl.Tail != nil {
		t.Error("dl.Head != nil || dl.Tail != nil")
	}
	if dl.Delete(el) != false {
		t.Error("empty list delete")
	}

	e1 := dl.AppendVal(1)
	e2 := dl.AppendVal(2)

	if !dl.Delete(e2) {
		t.Error("2nd el delet unsuccesfull")
	}
	e2 = dl.AppendVal(2)
	if !dl.Delete(dl.Tail) {
		t.Error("Tail delete unsuccesfull")
	}
	e2 = dl.AppendVal(2)

	if !dl.Delete(e1) {
		t.Error("1st el delet unsuccesfull")
	}
	if !dl.Delete(e2) {
		t.Error("2nd el delet unsuccesfull")
	}
	if dl.Head != nil || dl.Tail != nil {
		t.Error("dl.Head != nil || dl.Tail != nil")
	}
	if dl.Len != 0 {
		t.Error("dll.len != 0")
	}
	e1 = dl.AppendVal(1)
	e2 = dl.AppendVal(2)
	dl.AppendVal(3)
	if dl.Len != 3 {
		t.Error("dll.len != 3")
	}
	if !dl.Delete(e2) {
		t.Error("2nd el delet unsuccesfull")
	}
	if dl.Head.Next != dl.Tail || dl.Tail.Prev != dl.Head {
		t.Error("2nd el delet unsuccesfull")
	}
	dl.Clear()
	if dl.Len != 0 {
		t.Error("Clear()")
	}
	dl.Clear()
	if dl.Len != 0 {
		t.Error("Clear() of empty list")
	}
	dl.AppendVal(1)
	dl.AppendVal(2)
	el = dl.AppendVal(3)
	dl.AppendVal(4)
	dl.AppendVal(5)
	if dl.Len != 5 {
		t.Error("dll.len != 5")
	}
	dl.Delete(el)
	dl.Clear()
	dl.AppendVal(1)
	dl.AppendVal(2)
	dl.AppendVal(3)
	dl.AppendVal(4)
	dl.AppendVal(5)

	for j := 0; j < 10; j++ {
		iterator := dl.Iter()
		idx := 1
		for el, ok := iterator(); ok; el, ok = iterator() {
			if *&el.Value != idx {
				t.Error("iterator error")
			}
			idx++
		}
	}

}