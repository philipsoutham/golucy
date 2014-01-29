package golucy

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
	"time"
)

func TestIdxCreation(t *testing.T) {
	Convey("Index should do what it's told", t, func() {
		var index *Index
		defer index.Close()
		schema := NewSchema()
		defer schema.Close()
		So(schema, ShouldNotBeNil)
		index = makeIndex(schema, true, true, true)
		So(index.Path, ShouldNotBeBlank)
		So(index.Create, ShouldBeTrue)
		So(index.Truncate, ShouldBeTrue)
		So(index.Schema, ShouldNotBeNil)
		index = makeIndex(schema, true, false, true)
		So(index.Create, ShouldBeFalse)
		So(index.Truncate, ShouldBeTrue)
		index = makeIndex(schema, true, true, false)
		So(index.Create, ShouldBeTrue)
		So(index.Truncate, ShouldBeFalse)
		index = makeIndex(schema, true, false, false)
		So(index.Create, ShouldBeFalse)
		So(index.Truncate, ShouldBeFalse)
		//So()
	})

	Convey("Path should be invalid", t, func() {
		schema := NewSchema()
		defer schema.Close()
		index := makeIndex(schema, false, true, true)
		defer index.Close()
		So(index.Path, ShouldBeBlank)
	})

}

func TestSchema(t *testing.T) {
	Convey("Schema does what we want", t, func() {
		schema := NewSchema()
		defer schema.Close()
		index := makeIndex(schema, true, true, true)
		defer index.Close()
		index.NewIndexWriter()
		schema.AddField(NewIdField("id"))
		schema.AddField(NewFTField("content", "en"))
		So(len(schema.Fields), ShouldEqual, 2)
		So(schema.Fields[0].IndexType, ShouldEqual, StringType)
		So(schema.Fields[1].IndexType, ShouldEqual, FullTextType)
	})
}

func makeIndex(schema *Schema, valid, create, truncate bool) *Index {
	var tmpDir string

	if !valid {
		tmpDir = "/foo"
	}

	ixLocation, _ := ioutil.TempDir(tmpDir, fmt.Sprintf("lucy_test_%d", time.Now().UTC().UnixNano()))
	return NewIndex(ixLocation, create, truncate, schema)
}
