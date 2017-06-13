package structer

// Structer
// Mark tests
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "testing"

//import "log"

func TestMark(t *testing.T) {
	//var x int64 = -5
	//log.Print("--------------------------", int32(x))

	ords := map[string]int{"a": 1, "b": 2}
	//outList := []int{91, 92, 93}
	m := newMark(ords)
	m.addId(91, map[string]int{"a": 1, "b": 2})
	m.addId(94, map[string]int{"a": 6, "b": 5})
	m.addId(95, map[string]int{"a": 3, "b": 4})
	//arr, num := m.getCross(outList)
	//log.Print("--------------------------", arr, " ", num)
	//log.Print("--------------------------", m.lists)
}

func TestCountListsFields(t *testing.T) {
	ords := map[string]int{"a": 1, "b": 2}
	st := newMark(ords)

	if len(ords) != len(st.lists)-1 {
		t.Error("Error in the number of created fields in the variable `st.list`")
	}
}

func TestZeroCountListsFields(t *testing.T) {
	ords := map[string]int{}
	st := newMark(ords)

	if 1 != len(st.lists) {
		t.Error("No default field was created in `st.list`")
	}
}

/*
func TestAddIdCounter(t *testing.T) {
	st := newMark(map[string]int{})
	st.addId(91)
	st.addId(94)
	if st.count != 2 {
		t.Error("Invalid data in the added number counter`st.count`")
	}
}


func TestAddIdLenArr(t *testing.T) {
	st := newMark(map[string]int{})
	st.addId(91)
	st.addId(94)
	if len(st.arr) != st.count {
		t.Error("The number and number of numbers added do not match")
	}
}

func TestAddUnsafeCounter(t *testing.T) {
	st := newMark(map[string]int{})
	st.addUnsafe(91)
	st.addUnsafe(94)
	if st.count != 2 {
		t.Error("Invalid data in the unsafe added number counter`st.count`")
	}
}

func TestAddUnsafeLenArr(t *testing.T) {
	st := newMark(map[string]int{})
	st.addUnsafe(91)
	st.addUnsafe(94)
	if len(st.arr) != st.count {
		t.Error("The number and number of numbers unsafe added do not match")
	}
}

func TestDelIdCounter(t *testing.T) {
	st := newMark(map[string]int{})
	st.addId(91)
	st.addId(94)
	st.delId(94)
	if st.count != 1 {
		t.Error("Invalid data in the added and deleted number counter`st.count`")
	}
}

func TestDelIdLenArr(t *testing.T) {
	st := newMark(map[string]int{})
	st.addId(91)
	st.addId(94)
	st.delId(94)
	if len(st.arr) != st.count {
		t.Error("The number and number of numbers added and deleted do not match")
	}
}

func TestCrossOnlyAdd(t *testing.T) {
	outList := []int{91, 92, 93}
	st := newMark(map[string]int{})
	st.addId(91)
	st.addId(94)
	st.addId(95)
	arr, num := st.getCross(outList)
	if num != 1 {
		t.Error("Invalid number in the intersection of two sets (only Add)")
	}
	if len(arr) != 3 {
		t.Error("Invalid number in set (only Add)")
	}
}

func TestCrossAddDel(t *testing.T) {
	outList := []int{91, 92, 93}
	st := newMark(map[string]int{})
	st.addId(91)
	st.addId(92)
	st.addId(95)
	st.delId(92)
	arr, num := st.getCross(outList)
	if num != 1 {
		t.Error("Invalid number in the intersection of two sets (Add & Del)")
	}
	if len(arr) != 3 {
		t.Error("Invalid number in set (Add & Del)")
	}
}

func TestCrossZeroSet(t *testing.T) {
	outList := []int{91, 92, 93}
	st := newMark(map[string]int{})
	_, num := st.getCross(outList)
	if num != 0 {
		t.Error("Invalid number in the intersection of two sets (Add & Del)")
	}
}
*/
/*
func TestGetOrderedList(t *testing.T) {
	st := newMark(map[string]int{"a": 1, "b": 2})
	st.addId(91)
	st.addId(92)
	st.addId(95)
	arr := st.getOrderedList("a")
	t.Error(st.lists)
	if len(arr) != 3 {

	}
}
*/
