package structer

// Structer
// Query
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "log"

// newQuery - create a new Query-struct
func newQuery(db *EmbeDb) *Query {
	return &Query{
		db:     db,
		fields: make([]string, 0, RESERVED_SIZE_FOR_TAGS),
		sort:   "",
		asc:    0,
		from:   0,
		how:    HOW_MANY_STRUCT_RETURN,
	}
}

type Query struct {
	db     *EmbeDb
	fields []string
	sort   string
	asc    int
	from   int
	how    int
}

// ByFields - list of tags to be searched
func (q *Query) ByFields(where []string) *Query {
	q.fields = where
	return q
}

// OrderBy - tag (integer), by which the result will be sorted
func (q *Query) OrderBy(tag string, asc int) *Query {
	q.sort = tag
	q.asc = asc
	return q
}

// Limit - From the received list to return only ...
func (q *Query) Limit(from int, how int) *Query {
	q.from = from
	q.how = how
	return q
}

// Do - execute request (mandatory option)
func (q *Query) Do() []interface{} {
	return q.db.selectDo(q)
}
