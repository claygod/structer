package structer

// Structer
// Config
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

const (
	ASC = iota
	DESC
)

const RESERVED_SIZE_SLICE int = 100

// Query defailt
const (
	RESERVED_SIZE_FOR_TAGS int = 10
	HOW_MANY_STRUCT_RETURN int = 20
)

// For spec
const (
	TYPE_INT = iota
	TYPE_INT32
	TYPE_INT64
	TYPE_UINT
	TYPE_UINT32
	TYPE_UINT64
	TYPE_FLOAT32
	TYPE_FLOAT64
)
