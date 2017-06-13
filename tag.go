package structer

// Structer
// Tags
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "sync"
import "sort"
import "errors"
import "fmt"

//import "log"

// newTags - create a new Tags-struct
func newTags(sortIndexes map[string]int) *Tags {
	t := &Tags{
		subTags:     make(map[string]*Mark),
		mapTagValue: sortIndexes,
	}
	return t
}

// Tags - store subtags
// Хранилище всех тегов (под каждый тег выделяется Mark)
type Tags struct {
	sync.Mutex
	subTags     map[string]*Mark
	mapTagValue map[string]int // тут нужны только ключи
}

func (t *Tags) getTagList() []string {
	out := make([]string, len(t.subTags))
	t.Lock()
	i := 0
	for k := range t.subTags {
		out[i] = k
		i++
	}
	t.Unlock()
	return out
}

func (t *Tags) selectByTags(tagsNames []string, sortKey string) []int {
	if len(tagsNames) == 0 {
		return nil
	} else if len(tagsNames) == 1 {
		return t.subTags[tagsNames[0]].lists[""]
	}
	tagsCounts2 := make([]int, len(tagsNames))
	for i, tag := range tagsNames {
		if st, ok := t.subTags[tag]; ok {
			tagsCounts2[i] = int(uint64(len(st.arr))<<32) + i
		} else { // тут можно сделать выход по ошибке!
			// TO DO
		}
	}
	sort.Sort(sort.IntSlice(tagsCounts2))
	tagsOrdered := make([]string, len(tagsCounts2))

	for i, tag := range tagsCounts2 {
		n := int(uint32(tag))
		tagsOrdered[i] = tagsNames[n]
	}
	outList := t.subTags[tagsOrdered[0]].getOrderedList(sortKey)
	//if arr, ok := t.subTags[tagsOrdered[0]].lists[sortKey]; ok {
	//	outList = arr
	//} else {
	//	outList = t.subTags[tagsOrdered[0]].lists[""]
	//}
	cnt := len(outList)

	for i := 1; i < len(tagsOrdered); i++ {
		outList, cnt = t.subTags[tagsOrdered[i]].getCross(outList)
		if cnt == 0 {
			break
		}
	}
	return outList[:cnt]
}

func (t *Tags) addToTags(tagsNames []string, id int, mapTagValue map[string]int) bool {
	//log.Print("===============", sortIndexes)
	t.Lock()
	for _, tag := range tagsNames {
		if _, ok := t.subTags[tag]; !ok {
			t.subTags[tag] = newMark(t.mapTagValue)
		}
	}
	t.Unlock()
	for _, tag := range tagsNames {
		t.subTags[tag].addId(id, mapTagValue)
	}

	return true
}

func (t *Tags) addToTagsUnsafe(tagsNames []string, id int, mapTagValue map[string]int) bool {
	for _, tag := range tagsNames {
		if _, ok := t.subTags[tag]; !ok {
			t.subTags[tag] = newMark(t.mapTagValue)
		}

	}
	for _, tag := range tagsNames {
		t.subTags[tag].addUnsafe(id, mapTagValue)
	}

	return true
}

func (t *Tags) delFromTags(tagList []string, id int) {
	for _, tag := range tagList {
		ln := t.subTags[tag].delId(id)
		if ln == 0 {
			t.Lock()
			if _, ok := t.subTags[tag]; ok && t.subTags[tag].count == 0 {
				delete(t.subTags, tag)
			}
			t.Unlock()
		}
	}
}

func (t *Tags) tagsExists(tags []string) error {
	for _, tag := range tags {
		if _, ok := t.subTags[tag]; !ok {
			return errors.New(fmt.Sprintf("The tag `%s` does not exist", tag))
		}
	}
	return nil
}
