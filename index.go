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
package golucy

/*
#include "Clownfish/CharBuf.h"
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

type IndexType uint

type IndexOpenFlags int

type Index struct {
	Schema      *Schema
	Location    string
	Create      bool
	Truncate    bool
	lucyIndexer *C.LucyIndexer
}

type IndexValueColumn struct {
	Field string
	Value string
}

type IndexValueRow struct {
	Columns []*IndexValueColumn
}

type IndexWriter func(...*IndexValueColumn)

const (
	BlobType IndexType = iota
	FullTextType
	StringType
)

const IndexOpen IndexOpenFlags = 0
const (
	IndexCreate IndexOpenFlags = 1 << iota
	IndexTruncate
)

func NewIndex(schema *Schema, location string, create, truncate bool) *Index {
	return &Index{
		Schema:   schema,
		Location: location,
		Create:   create,
		Truncate: truncate,
	}
}

func (index *Index) GetIndexWriter(flags IndexOpenFlags) IndexWriter {
	idxLocation := cb_newf(index.Location)
	index.lucyIndexer = C.LucyIndexerNew(index.Schema.lucySchema, idxLocation, nil, C.int32_t(flags))
	C.DECREF(idxLocation)
	return func(ixValueColumns ...*IndexValueColumn) {
		doc := C.LucyDocNew(nil, 0)
		for _, ixColumn := range ixValueColumns {
			fieldName := cb_newf(ixColumn.Field)
			fieldValue := cb_new_from_utf8(ixColumn.Value)
			C.LucyDocStore(doc, fieldName, fieldValue)
			C.DECREF(fieldName)
			C.DECREF(fieldValue)
		}
		C.LucyIndexerAddDoc(index.lucyIndexer, doc, 1.0)
		C.DECREF(doc)
	}
}

// Commit any changes made to the index. Until this is called, none of the changes made during an indexing session are permanent.
// Calling commit() invalidates the Indexer, so if you want to make more changes you'll need a new one.
func (index *Index) Commit() {
	C.LucyIndexerCommit(index.lucyIndexer)
	C.DECREF(index.lucyIndexer)
}

func (index *Index) Close() {
	// if index.lucyIndexer != nil {
	// 	// need to figure out a way to ensure index.lucyIndexer is handled
	// 	// What is the value if it's already had DECREF called on it?
	// 	//C.DECREF(index.lucyIndexer)
	// }
}

// Optimize the index for search-time performance. This may take a while, as it
// can involve rewriting large amounts of data.
func (index *Index) Optimize() {
	C.LucyIndexerOptimize(index.lucyIndexer)
}

// Mark documents which contain the supplied term as deleted, so that they will be
// excluded from search results and eventually removed altogether. The change is not apparent to search apps until after commit() succeeds.
func (index *Index) DeleteByTerm(field, term string) {
	C.LucyIndexerDeleteByTerm(index.lucyIndexer, cb_newf(field), cb_new_from_utf8(term))
}

// Mark documents which match the supplied Query as deleted.
func (index *Index) DeleteByQuery(query string) {
	panic("not implemented")
}

func (index *Index) DeleteByDocId(docId int32) {
	C.LucyIndexerDeleteByDocId(index.lucyIndexer, C.int32_t(docId))
}
