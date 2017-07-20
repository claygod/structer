package extractor

// Extractor
// API
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "reflect"

// import "unsafe"
//import "fmt"

//import "log"

// NewExtractor - create a new Extractor-struct
// item - здесь структура как образец и т.д. (пригодится при проведении ревизии в движке)
func NewExtractor(item interface{}) (*Extractor, error) {
	s := &Extractor{
		structFields: make(map[string]reflect.StructField),
	}
	s.parseStruct(item)
	s.structType = reflect.TypeOf(item)

	return s, nil
}

// Extractor - спецификация
type Extractor struct {
	structType   reflect.Type
	structFields map[string]reflect.StructField
}

func (e *Extractor) Extract() *Kit {
	return &Kit{
		ext: e,
		arr: make([]*Query, 0),
	}
}

// Kit - спецификация
type Kit struct {
	ext  *Extractor
	item interface{}
	arr  []*Query
}

func Req() *Kit {
	return &Kit{arr: make([]*Query, 0)}
}
func (k *Kit) Of(r *Query) *Kit {
	k.arr = append(k.arr, r)
	return k
}
func (k *Kit) From(item interface{}) *Kit {
	k.item = item
	return k
}

func (k *Kit) Do() (map[string]interface{}, error) {
	return k.ext.extractData(k.item, k.arr)
}

// Query - спецификация
type Query struct {
	fieldName string
	fieldType reflect.Type
}

func Field(name string) *Query {
	return &Query{fieldName: name}
}
func (r *Query) Int() *Query {
	r.fieldType = reflect.TypeOf(int(1))
	return r
}
func (r *Query) String() *Query {
	r.fieldType = reflect.TypeOf(string(""))
	return r
}
func (r *Query) SlaceString() *Query {
	r.fieldType = reflect.TypeOf([]string{})
	return r
}
