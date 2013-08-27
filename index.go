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

#include "Lucy/Document/Doc.h"
#include "Lucy/Index/Indexer.h"

#define DECREF       cfish_Obj_decref

#define LucyIndexer lucy_Indexer
#define LucyIndexerNew lucy_Indexer_new
#define LucyIndexerCREATE lucy_Indexer_CREATE
#define LucyIndexerTRUNCATE lucy_Indexer_TRUNCATE
#define LucyDocStore LUCY_Doc_Store
#define LucyDocNew lucy_Doc_new
#define LucyIndexerAddDoc LUCY_Indexer_Add_Doc
#define LucyIndexerCommit LUCY_Indexer_Commit
*/
import "C"

type IndexType uint

type IndexOpenFlags int

type Index struct {
	Schema      *Schema
	Location    string
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
