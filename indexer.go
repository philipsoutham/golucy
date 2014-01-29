package golucy

import (
	"runtime"
)

// Copyright 2013 Philip Southam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
#include "Clownfish/Obj.h"
#define DECREF       cfish_Obj_decref

#include "Lucy/Document/Doc.h"
#define LucyDocStore LUCY_Doc_Store
#define LucyDocNew lucy_Doc_new

#include "Lucy/Index/Indexer.h"
#define LucyIndexer lucy_Indexer
#define LucyIndexerNew lucy_Indexer_new
#define LucyIndexerOptimize LUCY_Indexer_Optimize
#define LucyIndexerDeleteByTerm LUCY_Indexer_Delete_By_Term
#define LucyIndexerDeleteByQuery LUCY_Indexer_Delete_By_Query
#define LucyIndexerDeleteByDocId LUCY_Indexer_Delete_By_Doc_ID
#define LucyIndexerCREATE lucy_Indexer_CREATE
#define LucyIndexerTRUNCATE lucy_Indexer_TRUNCATE
#define LucyIndexerAddDoc LUCY_Indexer_Add_Doc
#define LucyIndexerCommit LUCY_Indexer_Commit

#include "Lucy/Search/Searcher.h"
#define LucySearcherGleanQuery LUCY_Searcher_Glean_Query

*/
import "C"

type Index struct {
	Path     string
	Create   bool
	Truncate bool
	Schema   *Schema
}

type IndexWriter struct {
	Index       *Index
	lucyIndexer *C.LucyIndexer
}

const (
	indexOpen   uint8 = 0
	indexCreate uint8 = 1 << iota
	indexTruncate
)

// best for creating
func NewIndex(path string, create, truncate bool, schema *Schema) *Index {
	return &Index{
		Path:     path,
		Create:   create,
		Truncate: truncate,
		Schema:   schema,
	}
}

// used just for reading
func OpenIndex(path string) *Index {
	return &Index{Path: path}
}

func (index *Index) Close() {
	// need to make sure everthing has been torn down
}

// Preferred method of creating an `IndexWriter`
// Set's up all the necessary C bindings.
func (index *Index) NewIndexWriter() *IndexWriter {
	flags := indexOpen
	if index.Create {
		flags |= indexCreate
	}
	if index.Truncate {
		flags |= indexTruncate
	}
	ixLocation := cb_news(index.Path)
	defer C.DECREF(ixLocation)
	ixWriter := &IndexWriter{
		Index:       index,
		lucyIndexer: C.LucyIndexerNew(index.Schema.lucySchema, ixLocation, nil, C.int32_t(flags)),
	}
	runtime.SetFinalizer(ixWriter, freeIndexWriter)
	return ixWriter
}

func (ixWriter *IndexWriter) AddDoc(doc Document) {
	lDoc := C.LucyDocNew(nil, 0) // Are these sane defaults?
	for k, v := range doc {
		name := cb_news(k)
		value := cb_new_from_utf8(v)
		C.LucyDocStore(lDoc, name, value)
		C.DECREF(name)
		C.DECREF(value)
	}
	C.LucyIndexerAddDoc(ixWriter.lucyIndexer, lDoc, 1.0) // Is 1.0 a sane default?
	C.DECREF(lDoc)
}

func (ixWriter *IndexWriter) AddDocs(docs ...Document) { // should this be []Document or ...Document?
	for _, doc := range docs {
		ixWriter.AddDoc(doc)
	}
}

func (ixWriter *IndexWriter) Commit() {
	C.LucyIndexerCommit(ixWriter.lucyIndexer)
}

func (ixWriter *IndexWriter) Close() {
	// Should this be here or in Commit?
	if ixWriter.lucyIndexer != nil {
		C.DECREF(ixWriter.lucyIndexer)
		ixWriter.lucyIndexer = nil
	}
}

func freeIndexWriter(ixWriter *IndexWriter) {
	ixWriter.Close()

}
