package structer

// Structer
// Api
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "log"
import "errors"
import "encoding/gob"
import "os"
import "sync"
import "fmt"

// New - create a new Structer
// item - это образец хранящейся структуры
// id - имя поля, которое будет id это поле должно быть строковым)
// tags - список тегов, это должны быть только  int, string, []string поля структуры
// из по int полям можно будет потом сортировать, а по string и []string делать
// запросы на выборки.
func New(item interface{}, id string, tags []string) (*Structer, error) {
	spec, err := newSpec(item, id, tags)
	if err != nil {
		return nil, err
	}
	mapTagValue := spec.getSortTags(item)
	e := &Structer{
		tags:    newTags(mapTagValue),
		spec:    spec,
		storage: newStorage(),
		index:   newIndex(),
	}
	return e, nil
}

// Structer - a structures storage
type Structer struct {
	sync.Mutex
	tags    *Tags
	spec    *Spec
	storage *Storage
	index   *Index
}

// Add - add structure to storage
// This function is unsafe, since the user "outside"
// can change or delete the structure,
// and the database will not know anything about it.
func (e *Structer) Add(item interface{}) error {
	mapTagValue := e.spec.getSortTags(item)
	// проверка на размер значений - не больше 32-битного int
	if err := e.checkSortTagSize(mapTagValue); err != nil {
		return err
	}
	id := e.spec.getId(item)
	k := e.storage.addItem(item)
	e.index.addId(id, k)
	tgs := e.spec.getTags(item)
	// mapTagValue := e.spec.getMapTagValue(item)
	e.tags.addToTags(tgs, k, mapTagValue)
	return nil
}

// AddUnsafe - add structure to storage (unsafe)
func (e *Structer) AddUnsafe(item interface{}) error {
	id := e.spec.getId(item)
	k := e.storage.addItemUnsafe(item)
	e.index.addIdUnsafe(id, k)
	tgs := e.spec.getTags(item)
	mapTagValue := e.spec.getSortTags(item)
	e.tags.addToTagsUnsafe(tgs, k, mapTagValue)
	return nil
}

// Update - replace structure
func (e *Structer) Update(itemNew interface{}) error {
	id := e.spec.getId(itemNew)
	num := e.index.getNumForId(id)
	if num < 0 {
		return errors.New("The record you want to transfer does not exist!")
	}
	itemOld := e.storage.getItem(num)
	tagsOld := e.spec.getTags(itemOld)
	tagsNew := e.spec.getTags(itemNew)
	// из списка делаем массив
	arrNew := make(map[string]bool)
	for _, tag := range tagsNew {
		arrNew[tag] = true
	}
	// формируем список на удаление
	listDel := make([]string, 0)
	for _, tag := range tagsOld {
		if _, ok := arrNew[tag]; ok {
			delete(arrNew, tag)
		} else {
			listDel = append(listDel, tag)
		}
	}
	// формируем список новых тегов
	listAdd := make([]string, 0)
	for tag := range arrNew {
		listAdd = append(listAdd, tag)
	}
	// удаляем старые теги
	e.tags.delFromTags(listDel, num)

	// меняем item
	e.storage.updateItem(itemNew, num)

	// добавляем новые теги
	mapTagValue := e.spec.getSortTags(itemNew)
	e.tags.addToTags(listAdd, num, mapTagValue)

	return nil
}

// Del - delete structure from storage
func (e *Structer) Del(id string) error {
	num := e.index.getNumForId(id)
	item := e.storage.getItem(num)
	if item == nil {
		return errors.New("There is no such record!")
	}
	tags := e.spec.getTags(item)
	e.tags.delFromTags(tags, num)
	e.index.delId(id)
	return nil
}

// Get - find structure by identifier
func (e *Structer) Get(id string) interface{} {
	num := e.index.getNumForId(id)
	return e.storage.getItem(num)
}

// Save - save the data store in a file (gob)
func (e *Structer) Save(path string) error {
	if e.fileExists(path) {
		os.Remove(path)
	}
	// open output file
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	encoder := gob.NewEncoder(fo)

	for _, v := range e.storage.arr {
		if _, ok := e.index.arr[e.spec.getId(v)]; ok {
			err := encoder.Encode(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Revision - garbage collection
func (e *Structer) Revision() error {
	newDb, err := New(
		e.spec.item,
		e.spec.idName,
		e.spec.sourceTags,
	)
	if err != nil {
		return err
	}

	e.Lock()
	//fmt.Print("\r\n В старой базе): ", len(e.storage.arr))
	for _, v := range e.storage.arr {
		if _, ok := e.index.arr[e.spec.getId(v)]; ok {
			newDb.AddUnsafe(v)
		}
	}
	//fmt.Print("\r\n В новой базе): ", len(newDb.storage.arr))
	e.index = newDb.index
	e.spec = newDb.spec
	e.storage = newDb.storage
	e.tags = newDb.tags
	e.Unlock()

	//e = newDb

	return nil
}

// Find - create a new search query (by tags)
func (e *Structer) Find() *Query {
	return newQuery(e)
}

// Count - number of stored structures
func (e *Structer) Count() int {
	cnt, _ := e.index.countAndDeleted()
	return cnt
}

// Fragmentation - percentage of fragmentation, the ratio of the unfragmented amount to the total.
func (e *Structer) Fragmentation() int {
	cnt, dltd := e.index.countAndDeleted()
	//fmt.Print("\r\n ==== ", cnt, dltd)
	return dltd * 100 / (cnt + dltd)
}

func (e *Structer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (e *Structer) selectDo(q *Query) ([]interface{}, error) {
	// Request validation
	//  len tags
	if len(q.fields) == 0 {
		return make([]interface{}, 0), errors.New("Tag list is empty")
	}
	// tag exists
	if err := e.tags.tagsExists(q.fields); err != nil {
		return make([]interface{}, 0), err
	}
	// limits
	if q.from < 0 || q.from >= len(e.index.arr) || q.how < 0 {
		return make([]interface{}, 0), errors.New("Error in request limits (from, how)")
	}
	// sort by
	if _, ok := e.tags.mapTagValue[q.sort]; !ok {
		return make([]interface{}, 0), errors.New(fmt.Sprintf("By tag `%s` you can not sort", q.sort))
	}
	return e.storage.listItems(e.limitIds(e.tags.selectByTags(q.fields, q.sort), q.from, q.how, q.asc), q.asc), nil

}

func (e *Structer) limitIds(tags []int, from int, how int, asc int) []int {
	ln := len(tags)
	if how < 1 || from < 0 || from >= ln { //
		return []int{}
	}

	if asc == ASC {
		to := from + how
		if to > ln {
			to = ln
		}
		if from > ln {
			from = ln
		}
		return tags[from:to]
	}
	from2 := ln - from - how
	if from2 > ln {
		from2 = ln
	} else if from2 < 0 {
		from2 = 0
	}
	to := ln - from
	if to > ln {
		to = ln
	}
	return tags[from2:to]
}

/*
func (e *Structer) limitItems(items []interface{}, from int, how int) []interface{} {
	if how > len(items) {
		how = len(items)
	}
	if from > len(items) {
		from = len(items)
	}
	return items[from:how]
}
*/

// checkSortTagSize - проверка размера int, не должен быть больше 32 бит
func (e *Structer) checkSortTagSize(mapTagValue map[string]int) error {
	for k, v := range mapTagValue {
		if v > 2147483647 {
			return errors.New(fmt.Sprintf("The value of the field `%s` for sorting exceeds 2147483647", k))
		}
	}
	return nil
}
