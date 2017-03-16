package structer

// Structer
// Index
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "sync"

// newIndex - create a new Index-struct
func newIndex() *Index {
	ix := &Index{arr: make(map[string]int)}
	return ix
}

type Index struct {
	sync.Mutex
	arr map[string]int
}

func (ix *Index) addId(id string, num int) bool {
	ix.Lock()
	if _, ok := ix.arr[id]; ok {
		ix.Unlock()
		return false
	}
	ix.arr[id] = num
	ix.Unlock()
	return true
}

func (ix *Index) addIdUnsafe(id string, num int) bool {
	if _, ok := ix.arr[id]; ok {
		ix.Unlock()
		return false
	}
	ix.arr[id] = num
	return true
}

func (ix *Index) delId(id string) bool {
	ix.Lock()
	if _, ok := ix.arr[id]; !ok {
		ix.Unlock()
		return false
	}
	delete(ix.arr, id)
	ix.Unlock()
	return true
}

func (ix *Index) getNumForId(id string) int {
	ix.Lock()
	if _, ok := ix.arr[id]; !ok {
		ix.Unlock()
		return -1
	}
	ix.Unlock()
	return ix.arr[id]
}