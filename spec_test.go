package structer

// Structer
// Spec tests
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "testing"

func TestSpecCorrect(t *testing.T) {
	item := SpecStructTest1{
		Id:   "abc",
		Pub:  true,
		Date: 123,
		Text: "Acme",
		Tags: []string{"Label1", "Label2"},
	}
	spec, err := newSpec(item, "abc", []string{"Date", "Tags"})
	if err != nil {
		t.Error(err)
	}
	if spec.getId(item) != "abc" {
		t.Error("The identifier is incorrectly extracted from the structure")
	}
	tgs := spec.getTags(item)
	if len(tgs) != 2 {
		t.Error("Error in the number of tags (select)")
	}
	if tgs[0] != "TagsLabel1" || tgs[1] != "TagsLabel2" {
		t.Error("Error in tag name")
	}
	stgs := spec.getSortTags(item)
	if len(stgs) != 1 {
		t.Error("Error in the number of tags (sort)")
	}
	if par, ok := stgs["Date"]; ok {
		if par != 123 {
			t.Error("Invalid field value")
		}
	} else {
		t.Error("Missing field")
	}
}

type SpecStructTest1 struct {
	Id   string
	Pub  bool
	Date int
	Text string
	Tags []string
}
