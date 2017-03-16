package structer

// Structer
// Api
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "log"
import "errors"
import "encoding/gob"
import "os"
import "sync"

// New - create a new EmbeDb-struct
func New(item interface{}, id string, tags []string) (*EmbeDb, error) {
	spec, err := newSpec(item, id, tags)
	sortIndexes := spec.getSortIndexes(item)
	if err != nil {
		return nil, err
	}
	e := &EmbeDb{
		tags:    newTags(sortIndexes),
		spec:    spec,
		storage: newStorage(),
		index:   newIndex(),
	}
	return e, nil
}

// EmbeDb
type EmbeDb struct {
	sync.Mutex
	tags    *Tags
	spec    *Spec
	storage *Storage
	index   *Index
}

func (e *EmbeDb) Add(item interface{}) error {
	id := e.spec.getId(item)
	k := e.storage.addItem(item)
	e.index.addId(id, k)
	tgs := e.spec.getTags(item)
	e.tags.addToTags(tgs, k)
	return nil
}

func (e *EmbeDb) AddUnsafe(item interface{}) error {
	id := e.spec.getId(item)
	k := e.storage.addItemUnsafe(item)
	e.index.addIdUnsafe(id, k)
	tgs := e.spec.getTags(item)
	e.tags.addToTagsUnsafe(tgs, k)
	return nil
}

func (e *EmbeDb) Update(itemNew interface{}) error {
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
	for tag, _ := range arrNew {
		listAdd = append(listAdd, tag)
	}
	// удаляем старые теги
	e.tags.delFromTags(listDel, num)

	// меняем item
	e.storage.updateItem(itemNew, num)

	// добавляем новые теги
	e.tags.addToTags(listAdd, num)

	return nil
}

func (e *EmbeDb) Del(id string) error {
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

func (e *EmbeDb) Get(id string) interface{} {
	num := e.index.getNumForId(id)
	return e.storage.getItem(num)
}

func (e *EmbeDb) Save(path string) error {
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

func (e *EmbeDb) Revision() error {
	newDb, err := New(
		e.spec.item,
		e.spec.idName,
		e.spec.sourceTags,
	)
	if err != nil {
		return err
	}

	oldDb := e
	oldDb.Lock()
	for _, v := range e.storage.arr {
		if _, ok := e.index.arr[e.spec.getId(v)]; ok {
			newDb.AddUnsafe(v)
		}
	}
	oldDb.Unlock()
	e = newDb

	return nil
}

func (e *EmbeDb) Select() *Query {
	return newQuery(e)
}

func (e *EmbeDb) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (e *EmbeDb) selectDo(q *Query) []interface{} {
	if len(q.fields) == 0 {
		return make([]interface{}, 0)
	}
	return e.storage.listItems(e.limitIds(e.tags.selectByTags(q.fields, q.sort), q.from, q.how, q.asc), q.asc)

}

func (e *EmbeDb) limitIds(tags []int, from int, how int, asc int) []int {
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
	} else {

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

}

func (e *EmbeDb) limitItems(items []interface{}, from int, how int) []interface{} {
	if how > len(items) {
		how = len(items)
	}
	if from > len(items) {
		from = len(items)
	}
	return items[from:how]
}