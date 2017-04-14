package structer

// Structer
// Mark
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "sync"

//import "log"

func newMark(ords map[string]int) *Mark {
	//log.Print("Субтег получил список ords: ", ords)
	m := &Mark{
		arr:   make(map[int]bool),
		lists: make(map[string][]int),
		blum:  &Blum{},
	}
	for k := range ords {
		m.lists[k] = make([]int, 0, RESERVED_SIZE_SLICE)
	}
	m.lists[""] = make([]int, 0, RESERVED_SIZE_SLICE)
	return m
}

// Mark - Repository of sub-tags (sections)
type Mark struct {
	sync.Mutex
	count       int
	arr         map[int]bool
	lists       map[string][]int
	blum        *Blum
	sortIndexes map[string]int
}

func (m *Mark) addId(id int, sortIndexes map[string]int) bool {
	m.Lock()
	if _, ok := m.arr[id]; ok {
		m.Unlock()
		return false
	}
	m.arr[id] = true
	m.blum.addId(id)
	// в v значения по которым сортировать
	for k, v := range sortIndexes {
		v2 := int(uint64(uint32(int32(v)))<<32 + uint64(uint32(id)))
		if arr, ok := m.lists[k]; ok {
			if len(arr) == 0 {
				m.lists[k] = append(m.lists[k], v2)
				continue
			}
			flag := true
			ln := len(arr)
			average := arr[0] / 2
			average += arr[ln-1] / 2
			if v2 > average {
				for i := ln - 1; i < 0; i-- {
					if arr[i] < v2 {
						arr = append(arr[:i], append([]int{v2}, arr[i:]...)...)
						flag = false
						break
					}
				}
			} else {
				for x, y := range arr {
					if y < v2 {
						arr = append(arr[:x], append([]int{v2}, arr[x:]...)...)
						flag = false
						break
					}
				}
			}

			if flag {
				arr = append(arr, v2)
			}
			m.lists[k] = arr
		}
	}
	m.lists[""] = append(m.lists[""], id)

	m.count++
	m.Unlock()
	return true
}
func (m *Mark) addUnsafe(id int, sortIndexes map[string]int) bool {
	if _, ok := m.arr[id]; ok {
		return false
	}
	m.arr[id] = true
	m.blum.addId(id)
	// в v значения по которым сортировать
	for k, v := range sortIndexes {
		v2 := int(uint64(uint32(int32(v)))<<32 + uint64(uint32(id)))
		if arr, ok := m.lists[k]; ok {
			if len(arr) == 0 {
				m.lists[k] = append(m.lists[k], v2)
				continue
			}
			flag := true
			ln := len(arr)
			average := arr[0] / 2
			average += arr[ln-1] / 2
			if v2 > average {
				for i := ln - 1; i < 0; i-- {
					if arr[i] < v2 {
						arr = append(arr[:i], append([]int{v2}, arr[i:]...)...)
						flag = false
						break
					}
				}
			} else {
				for x, y := range arr {
					if y < v2 {
						arr = append(arr[:x], append([]int{v2}, arr[x:]...)...)
						flag = false
						break
					}
				}
			}

			if flag {
				arr = append(arr, v2)
			}
			m.lists[k] = arr
		}
	}
	m.lists[""] = append(m.lists[""], id)

	m.count++
	return true
}

func (m *Mark) delId(id int) int {
	m.Lock()
	if _, ok := m.arr[id]; ok {
		delete(m.arr, id)
		for tag, lst := range m.lists {
			for k, v := range lst {
				v2 := int(uint32(v))
				if v2 == id {
					m.lists[tag] = append(m.lists[tag][:k], m.lists[tag][k+1:]...)
					break
				}
			}
		}
	}
	m.count--
	ln := m.count
	m.Unlock()
	return ln
}

func (m *Mark) getCross(outList []int) ([]int, int) {
	cnt := 0
	var k int
	m.Lock()
	for u := 0; u < len(outList); u++ {
		k = int(uint32(outList[u])) // берём только правую часть цифры! она - номер, а левая - это значение по которому сортировали

		if m.blum.checkId(k) {
			if _, ok := m.arr[k]; ok {
				outList[cnt] = k
				cnt++
			}
		}
	}
	m.Unlock()
	return outList, cnt
}

func (m *Mark) getOrderedList(sortKey string) []int {
	if _, ok := m.lists[sortKey]; ok {
		return m.lists[sortKey]
	} else {
		return m.lists[""]
	}
}

/*
// Blum - hashes
type Blum [65536]uint8

func (b *Blum) addId(id int) {
	//uid := uint64(id)
	key := uint16(uint64(id))
	hash := b[key]
	shift := (uint64(id) << 45) >> 61
	in := uint8(1) << shift
	hash = hash | in
	b[key] = hash
}

func (b *Blum) checkId(id int) bool {
	key := uint16(uint64(id))
	hash := b[key]
	shift := (uint64(id) << 45) >> 61
	in := uint8(1) << shift
	if in == hash&in {
		return true
	}
	return false
}
*/
