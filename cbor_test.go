package main

import (
	"fmt"
	"testing"
)

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestDotPathToSlice(t *testing.T) {
	res := DotPathToSlice("my.data.path")
	expected := []string{"my", "data", "path"}
	if !Equal(res, expected) {
		t.Errorf("DotPathToSlice(\"my.data.path\") failed, expected %v, got %v", expected, res)
	}
}

func TestSimpleSearch(t *testing.T) {
	object := map[string]interface{}{"a": "AA", "b": "BB", "c": "CCC"}
	c := Container{object}
	var expected interface{}

	for k, v := range object {
		res, _ := c.Search(DotPathToSlice(k)...)
		expected = v
		if res.object != expected {
			t.Errorf("c.Search(DotPathToSlice(\"%s\")...) failed, expected %v, got %v", k, expected, res.object)
		}
	}
}

func TestSimpleSearch2(t *testing.T) {
	object := map[interface{}]interface{}{"a": "AA", "b": "BB", "c": "CCC"}
	c := Container{object}
	var expected interface{}

	for k, v := range object {
		res, _ := c.Search(DotPathToSlice(fmt.Sprintf("%v", k))...)
		expected = v
		if res.object != expected {
			t.Errorf("c.Search(DotPathToSlice(\"%s\")...) failed, expected %v, got %v", k, expected, res.object)
		}
	}
}

func TestSimpleSearchArray(t *testing.T) {
	object := []interface{}{"AAA", "BBB", "CCC"}
	c := Container{object}
	var expected interface{}

	for k, v := range object {
		res, _ := c.Search(DotPathToSlice(fmt.Sprintf("%v", k))...)
		expected = v
		if res.object != expected {
			t.Errorf("c.Search(DotPathToSlice(\"%v\")...) failed, expected %v, got %v", k, expected, res.object)
		}
	}
}

func TestNestedSearch(t *testing.T) {
	jsonData := []byte(`{"Name":"Eve","Age":6,"Parents":[{"Name":"Alice","Age":34},{"Name":"Bob","Age":36}]}`)
	var v interface{}
	json.Unmarshal(jsonData, &v)
	data := v.(map[string]interface{})
	c := Container{data}

	tests := map[string]string{
		"Name":           "Eve",
		"Age":            "6",
		"Parents.0.Name": "Alice",
		"Parents.0.Age":  "34",
		"Parents.1.Name": "Bob",
		"Parents.1.Age":  "36",
	}

	for path, expected := range tests {
		res, _ := c.Search(DotPathToSlice(path)...)
		if fmt.Sprintf("%v", res.object) != expected {
			t.Errorf("c.Search(DotPathToSlice(\"%v\")...) failed, expected %v, got %v", path, expected, res.object)
		}
	}
}

func TestNestedSet(t *testing.T) {
	jsonData := []byte(`{"Name":"Eve","Age":6,"Parents":[{"Name":"Alice","Age":34},{"Name":"Bob","Age":36}]}`)
	var v interface{}
	json.Unmarshal(jsonData, &v)
	data := v.(map[string]interface{})
	c := Container{data}

	tests := map[string]string{
		"Name":           "Cindy",
		"Age":            "8",
		"Parents.0.Name": "Andy",
		"Parents.0.Age":  "32",
		"Parents.1.Name": "Brenda",
		"Parents.1.Age":  "38",
	}

	for path, expected := range tests {
		c.Set(expected, DotPathToSlice(path)...)
		res, _ := c.Search(DotPathToSlice(path)...)
		if res.object != expected {
			t.Errorf("c.Set(\"%v\", DotPathToSlice(\"%v\")...) failed, expected %v, got %v", expected, path, expected, res.object)
		}
	}
}

func runDecodePathTests(t *testing.T, cbor string, tests map[string]string, isBase64 bool) {
	for path, expected := range tests {
		res, _ := DecodePath(cbor, path, isBase64)
		if fmt.Sprintf("%v", res) != expected {
			t.Errorf("UpdatePath() failed, expected %v, got %v", expected, res)
		}
	}
}

func TestDecodeSimplePathHex(t *testing.T) {
	cbor := "a56161614161626142616361436164614461656145"
	tests := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
		"d": "D",
		"e": "E",
	}
	runDecodePathTests(t, cbor, tests, false)
}

func TestDecodeSimplePathBase64(t *testing.T) {
	cbor := "pWFhYUFhYmFCYWNhQ2FkYURhZWFF"
	tests := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
		"d": "D",
		"e": "E",
	}
	runDecodePathTests(t, cbor, tests, true)
}

func TestDecodeNestedPathHex(t *testing.T) {
	cbor := "a26161016162820203"
	tests := map[string]string{
		"b.0": "2",
		"b.1": "3",
	}
	runDecodePathTests(t, cbor, tests, false)
}

func TestDecodeNestedPathBase64(t *testing.T) {
	cbor := "omFhAWFiggID"
	tests := map[string]string{
		"b.0": "2",
		"b.1": "3",
	}
	runDecodePathTests(t, cbor, tests, true)
}

func runUpdatePathTests(t *testing.T, cbor string, tests map[string]string, isBase64 bool) {
	for path, expected := range tests {
		res, _ := UpdatePath(cbor, path, expected, isBase64)
		res2, _ := DecodePath(res, path, isBase64)
		if res2 != expected {
			t.Errorf("UpdatePath() failed, expected %v, got %v", expected, res2)
		}
	}
}

func TestUpdatePathSimplePathHex(t *testing.T) {
	cbor := "a56161614161626142616361436164614461656145"
	tests := map[string]string{
		"a": "AA",
		"b": "BB",
		"c": "CC",
		"d": "DD",
		"e": "EE",
	}
	runUpdatePathTests(t, cbor, tests, false)
}

func TestUpdatePathSimplePathBase64(t *testing.T) {
	cbor := "pWFhYUFhYmFCYWNhQ2FkYURhZWFF"
	tests := map[string]string{
		"a": "AAA",
		"b": "BBB",
		"c": "CCC",
		"d": "DDD",
		"e": "EEE",
	}
	runUpdatePathTests(t, cbor, tests, true)
}

func TestUpdateNestedPathHex(t *testing.T) {
	cbor := "a26161016162820203"
	tests := map[string]string{
		"b.0": "changed",
		"b.1": "tothis",
	}
	runUpdatePathTests(t, cbor, tests, false)
}

func TestUpdateNestedPathBase64(t *testing.T) {
	cbor := "omFhAWFiggID"
	tests := map[string]string{
		"b.0": "changed",
		"b.1": "tothis",
	}
	runUpdatePathTests(t, cbor, tests, true)
}
