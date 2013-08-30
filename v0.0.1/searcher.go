package golucy

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

#include "Lucy/Search/IndexSearcher.h"
#define LucyIndexSearcher lucy_IndexSearcher
#define LucyIxSearcherNew lucy_IxSearcher_new
#define LucyIxSearcherHits LUCY_IxSearcher_Hits
#define LucyIxSearcherGetSchema LUCY_IxSearcher_Get_Schema

#include "Lucy/Analysis/EasyAnalyzer.h"
#define LucyEasyAnalyzerNew lucy_EasyAnalyzer_new

#include "Lucy/Plan/Schema.h"
#define LucySchema lucy_Schema
#define LucySchemaAllFields LUCY_Schema_All_Fields

#include "Lucy/Search/QueryParser.h"
#define LucyQueryParser lucy_QueryParser
#define LucyQParserNew lucy_QParser_new
#define LucyQParserParse LUCY_QParser_Parse

#include "Lucy/Document/HitDoc.h"
#define LucyHitDoc lucy_HitDoc
#define LucyHitDocExtract LUCY_HitDoc_Extract

#include "Lucy/Search/Hits.h"
#define LucyHitsTotal LUCY_Hits_Total_Hits
#define LucyHitsNext LUCY_Hits_Next

#include "Lucy/Search/Query.h"
#define LucyQuery lucy_Query

*/
import "C"

type Query struct {
	QueryStr   string
	lucySchema *C.LucySchema // we're now carrying this around in 2 places :-/
	lucyQuery  *C.LucyQuery
}

type IndexReader struct {
	Index        *Index
	lucySearcher *C.LucyIndexSearcher
}

func (index *Index) NewIndexReader() *IndexReader {
	ixLocation := cb_newf(index.Path)
	defer C.DECREF(ixLocation)
	return &IndexReader{Index: index, lucySearcher: C.LucyIxSearcherNew(ixLocation)}
}

func (ixReader *IndexReader) ParseQuery(queryStr string) *Query {
	lucySchema := C.LucyIxSearcherGetSchema(ixReader.lucySearcher)
	language := cb_newf("en") // should be configurable
	defer C.DECREF(language)
	analyzer := C.LucyEasyAnalyzerNew(language)
	defer C.DECREF(analyzer)
	qp := C.LucyQParserNew(
		lucySchema,
		analyzer,                          //should this be configurable?
		cb_newf("AND"),                    // should be configurable
		C.LucySchemaAllFields(lucySchema), // should be configurable
	)
	defer C.DECREF(qp)
	return &Query{
		QueryStr:   queryStr,
		lucySchema: lucySchema,
		lucyQuery:  C.LucyQParserParse(qp, cb_new_from_utf8(queryStr)),
	}
}

func (ixReader *IndexReader) Search(query *Query, offset, limit uint, field string) (uint, []string) {
	// Need to add `includeMatchedTerms bool` parameter. Then figure out a
	// way to extract the matched terms. Should probably have some sort
	// of `Results` object/iterator so that we don't have to specify
	// offset/limit and where I can attach matched terms to the result.
	getField := cb_newf(field) // total hack, need to return more than one field
	defer C.DECREF(getField)
	hits := C.LucyIxSearcherHits(ixReader.lucySearcher, query.lucyQuery, C.uint32_t(offset), C.uint32_t(limit), nil)
	defer C.DECREF(hits)
	totalNumHits := uint(C.LucyHitsTotal(hits))
	num2Return := minUInt(limit, totalNumHits)
	results := make([]string, num2Return)
	var hit *C.LucyHitDoc
	for i := uint(0); i < num2Return; i++ {
		hit = C.LucyHitsNext(hits)
		if hit == nil {
			break
		}
		value_cb := C.LucyHitDocExtract(hit, getField, nil) // do i need to free this, what does the nil do?
		value := cb_ptr2char(value_cb)                      // do i need to free this
		results[i] = C.GoString(value)
		C.DECREF(hit)
	}
	return num2Return, results
}

func (ixReader *IndexReader) Close() {
	C.DECREF(ixReader.lucySearcher)
}

func (query *Query) Close() {
	C.DECREF(query.lucySchema)
	C.DECREF(query.lucyQuery)
}
