package extractor

// Extractor
// Tools
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "reflect"
import "unsafe"
import "errors"

//import "log"

func (e *Extractor) parseStruct(item interface{}) error {
	t1 := reflect.TypeOf(item)
	for i := 0; i < t1.NumField(); i++ {
		e.structFields[t1.Field(i).Name] = t1.Field(i)
	}
	return nil
}

func (e *Extractor) extractData(item interface{}, q []*Query) (map[string]interface{}, error) {
	if reflect.TypeOf(item) != e.structType {
		return nil, errors.New("The obtained structure of a different type")
	}
	out := make(map[string]interface{})

	for _, v := range q {
		if s, ok := e.structFields[v.fieldName]; ok {
			if s.Type != v.fieldType {
				return nil, errors.New("The requested and returned types do not match")
			}
			//log.Print(s.Type, " ")
			switch s.Type {
			case reflect.TypeOf(int(1)):
				out[v.fieldName] = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + s.Offset))
			case reflect.TypeOf(string("")):
				out[v.fieldName] = *(*string)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + s.Offset))
			case reflect.TypeOf([]string{}):
				out[v.fieldName] = *(*[]string)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + s.Offset))
			default:
				return nil, errors.New("Unsupported type!")

			}
		} else {
			return nil, errors.New("The field with this name is missing")
		}
	}
	//log.Print("Получил: ", out)
	return out, nil
}

// field structure
type field struct {
	Name   string
	Type   reflect.Type
	Offset uintptr
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	bad    int32
	unused int32
	fun    [1]uintptr // variable sized
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

type _type struct {
	name   string
	bits   uint
	signed bool
}

type name struct {
	bytes *byte
}

type imethod struct {
	name nameOff
	ityp typeOff
}

type nameOff int32
type typeOff int32
type textOff int32
