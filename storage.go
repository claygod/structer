package structer

// Structer
// Storage
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "sync"

//import "log"

// NewStorage - create a new SubStore-struct
func newStorage() *Storage {
	s := &Storage{arr: make([]interface{}, 0, RESERVED_SIZE_SLICE)}
	return s
}

// SubStore - хранилище субструктур (секций)
type Storage struct {
	sync.Mutex
	arr []interface{}
}

func (s *Storage) addItem(item interface{}) int {
	s.Lock()
	num := len(s.arr)
	s.arr = append(s.arr, item)
	s.Unlock()
	return num
}

func (s *Storage) addItemUnsafe(item interface{}) int {
	num := len(s.arr)
	s.arr = append(s.arr, item)
	return num
}

func (s *Storage) updateItem(item interface{}, num int) bool {
	s.Lock()
	if len(s.arr) <= num {
		s.Unlock()
		return false
	}
	s.arr[num] = item
	s.Unlock()
	return true
}

func (s *Storage) getItem(id int) interface{} {
	s.Lock()
	if len(s.arr) < id {
		s.Unlock()
		return nil
	}
	s.Unlock()
	return s.arr[id]
}

func (s *Storage) listItems(m []int, asc int) []interface{} {
	out := make([]interface{}, len(m))
	s.Lock()
	if asc == ASC {
		for i, u := (len(m) - 1), 0; i >= 0; i, u = i-1, u+1 {
			out[u] = s.arr[m[i]]
		}

	} else {
		for i, id := range m {
			out[i] = s.arr[id]
		}
	}
	s.Unlock()
	return out
}
