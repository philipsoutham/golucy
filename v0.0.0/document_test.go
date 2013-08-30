package golucy_test

import (
	golucy "github.com/philipsoutham/golucy/v0.0.0"
	"testing"
)

func TestNewDoc(t *testing.T) {
	doc := golucy.NewDocument()
	doc.Add("foo", "bar")
	if val, ok := doc["foo"]; !ok || val != "bar" {
		t.Log("not what we expected", doc)
		t.FailNow()
	}

	doc = golucy.NewDocument()
	doc["foo"] = "bar"
	if val, ok := doc["foo"]; !ok || val != "bar" {
		t.Log("not what we expected", doc)
		t.FailNow()
	}

	doc = golucy.Document{"foo": "bar"}
	if val, ok := doc["foo"]; !ok || val != "bar" {
		t.Log("not what we expected", doc)
		t.FailNow()
	}

	fields := doc.GetFields()
	if len(fields) != 1 && cap(fields) != 1 {
		t.Log("not what we expected", len(fields), cap(fields), len(doc))
		t.FailNow()
	}
	if fields[0] != "foo" {
		t.Log("not what we expected", fields[0], fields)
		t.FailNow()
	}
}
