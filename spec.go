package structer

// Structer
// Spec
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "reflect"
import "unsafe"
import "errors"
import "fmt"

//import "log"

// newSpec - create a new Spec-struct
func newSpec(item interface{}, cId string, cTags []string) (*Spec, error) {
	s := &Spec{
		item:            item,
		idName:          cId,
		fields:          make(map[string]field),
		offsetSortType:  make(map[string]int),
		offsetSortPtr:   make(map[string]uintptr),
		offsetTagsRoot:  make(map[string]uintptr),
		offsetTagsSlice: make(map[string]uintptr),
		sourceTags:      cTags,
	}
	if err := s.parseFields(item, cTags); err != nil {
		return nil, err
	}
	// Устанавливаем смещение для Id
	s.offsetId = s.fields[cId].Offset
	// Устанавливаем смещения
	for _, t := range cTags {
		if f, ok := s.fields[t]; ok {
			switch f.Type {
			case reflect.TypeOf(""):
				s.offsetTagsRoot[t] = f.Offset
			case reflect.TypeOf([]string{}):
				s.offsetTagsSlice[t] = f.Offset
			case reflect.TypeOf(int(1)):
				s.offsetSortType[t] = TYPE_INT
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(int32(1)):
				s.offsetSortType[t] = TYPE_INT32
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(int64(1)):
				s.offsetSortType[t] = TYPE_INT64
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(uint(1)):
				s.offsetSortType[t] = TYPE_UINT
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(uint32(1)):
				s.offsetSortType[t] = TYPE_UINT32
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(uint64(1)):
				s.offsetSortType[t] = TYPE_UINT64
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(float32(1)):
				s.offsetSortType[t] = TYPE_FLOAT32
				s.offsetSortPtr[t] = f.Offset
			case reflect.TypeOf(float64(1)):
				s.offsetSortType[t] = TYPE_FLOAT64
				s.offsetSortPtr[t] = f.Offset
			default:
				return nil, errors.New(fmt.Sprintf(`"%s" Field can not be indexed`, t))
			}
		} else {
			return nil, errors.New(fmt.Sprintf(`Incorrect structure field name: "%s"`, t))
		}
	}
	return s, nil
}

// Spec - спецификация
type Spec struct {
	item            interface{}
	itemName        string
	itemType        reflect.Type
	idName          string
	offsetId        uintptr
	offsetSortType  map[string]int
	offsetSortPtr   map[string]uintptr
	offsetTagsRoot  map[string]uintptr
	offsetTagsSlice map[string]uintptr
	fields          map[string]field
	test            interface{}
	sourceTags      []string //полученный при создании список тегов
}

func (s *Spec) getId(item interface{}) string {
	return *(*string)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + s.offsetId))
}

func (s *Spec) getSortParam(tag string) (int, uintptr, bool) {
	if tp, ok := s.offsetSortType[tag]; ok {
		return tp, s.offsetSortPtr[tag], true
	}
	return 0, 0, false
}

func (s *Spec) getTags(item interface{}) []string {
	out := make([]string, 0, RESERVED_SIZE_SLICE)
	// Root tags
	for k, t := range s.offsetTagsRoot {
		out = append(out, k+*(*string)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + t)))
	}
	// Slice tags
	for k, t := range s.offsetTagsSlice {
		slc := *(*[]string)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + t))
		for _, t2 := range slc {
			out = append(out, k+t2)
		}
	}
	return out
}
func (s *Spec) getSortIndexes(item interface{}) map[string]int {
	// параметры для индексации под сортировку
	sortIndexes := make(map[string]int)
	for k, t := range s.offsetSortPtr {
		sortIndexes[k] = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&item)).data)) + t))
	}
	return sortIndexes
}

func (s *Spec) parseFields(item interface{}, fields []string) error {
	t1 := reflect.TypeOf(item)
	v1 := reflect.ValueOf(item)
	v1 = reflect.Indirect(v1)
	s.itemType = v1.Type()
	for _, key := range fields {
		if t2, ok := t1.FieldByName(key); ok {
			s.fields[t2.Name] = field{
				Name:   t2.Name,
				Type:   t2.Type,
				Offset: t2.Offset,
			}
		} else {
			return errors.New("This field in the structure is not")
		}

	}
	s.offsetId = s.fields["Id"].Offset
	return nil
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
