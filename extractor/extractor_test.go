package extractor

// Extractor
// Tests
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "testing"

//import "reflect"

func TestExtractor1(t *testing.T) {
	item := specStructTest1{
		Id:   "abc",
		Pub:  true,
		Date: 123,
		Text: "Acme",
		Tags: []string{"Label1", "Label2"},
	}
	spec, err := NewExtractor(item)
	if err != nil {
		t.Error(err)
	}

	re, err := spec.Extract().From(item).
		Of(Field("Date").Int()).
		Of(Field("Text").String()).
		Of(Field("Tags").SlaceString()).
		Do()

	if err != nil {
		t.Error(err)
	}
	if d, ok := re["Date"]; !ok {
		t.Error("The required field (Date) was not found")
	} else if d.(int) != 123 {
		t.Error("A response was received with an incorrect value, it was expected: 123")
	}

	if d, ok := re["Text"]; !ok {
		t.Error("The required field (Text) was not found")
	} else if d.(string) != "Acme" {
		t.Error("A response was received with an incorrect value, it was expected: Acme")
	}

	if d, ok := re["Tags"]; !ok {
		t.Error("The required field (Tags) was not found")
	} else if d.([]string)[0] != "Label1" {
		t.Error("A response was received with an incorrect value, it was expected: Label1")
	}

}

type specStructTest1 struct {
	Id   string
	Pub  bool
	Date int
	Text string
	Tags []string
}
