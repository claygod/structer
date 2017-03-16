package structer

// Structer
// Tags
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "sync"

//import "log"

func newSubTag(ords map[string]int) *SubTag {
	//log.Print("Субтег получил список ords: ", ords)
	s := &SubTag{
		arr:   make(map[int]bool),
		lists: make(map[string][]int),
		blum:  &Blum{},
	}
	for k := range ords {
		s.lists[k] = make([]int, 0, RESERVED_SIZE_SLICE)
	}
	s.lists[""] = make([]int, 0, RESERVED_SIZE_SLICE)
	return s
}

// SubTag - Repository of sub-tags (sections)
type SubTag struct {
	sync.Mutex
	count       int
	arr         map[int]bool
	lists       map[string][]int
	blum        *Blum
	sortIndexes map[string]int
}

func (s *SubTag) addId(id int) bool {
	s.Lock()
	if _, ok := s.arr[id]; ok {
		s.Unlock()
		return false
	}
	s.arr[id] = true
	s.blum.addId(id)
	// в v значения по которым сортировать
	for k, v := range s.sortIndexes {
		v2 := int(uint64(uint32(int32(v)))<<32 + uint64(uint32(id)))
		if arr, ok := s.lists[k]; ok {
			if len(arr) == 0 {
				s.lists[k] = append(s.lists[k], v2)
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
			s.lists[k] = arr
		}
	}
	s.lists[""] = append(s.lists[""], id)

	s.count++
	s.Unlock()
	return true
}
func (s *SubTag) addUnsafe(id int) bool {
	if _, ok := s.arr[id]; ok {
		return false
	}
	s.arr[id] = true
	s.blum.addId(id)
	// в v значения по которым сортировать
	for k, v := range s.sortIndexes {
		v2 := int(uint64(uint32(int32(v)))<<32 + uint64(uint32(id)))
		if arr, ok := s.lists[k]; ok {
			if len(arr) == 0 {
				s.lists[k] = append(s.lists[k], v2)
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
			s.lists[k] = arr
		}
	}
	s.lists[""] = append(s.lists[""], id)

	s.count++
	return true
}

func (s *SubTag) delId(id int) int {
	s.Lock()
	if _, ok := s.arr[id]; ok {
		delete(s.arr, id)
		for tag, lst := range s.lists {
			for k, v := range lst {
				v2 := int(uint32(v))
				if v2 == id {
					s.lists[tag] = append(s.lists[tag][:k], s.lists[tag][k+1:]...)
					break
				}
			}
		}
	}
	s.count--
	ln := s.count
	s.Unlock()
	return ln
}

func (s *SubTag) getCross(outList []int) ([]int, int) {
	cnt := 0
	var k int
	s.Lock()
	for u := 0; u < len(outList); u++ {
		k = int(uint32(outList[u])) // берём только правую часть цифры! она - номер, а левая - это значение по которому сортировали

		if s.blum.checkId(k) {
			if _, ok := s.arr[k]; ok {
				outList[cnt] = k
				cnt++
			}
		}
	}
	s.Unlock()
	return outList, cnt
}

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
