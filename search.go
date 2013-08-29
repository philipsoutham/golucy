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
#include <stdlib.h>
#include "Lucy/Search/Hits.h"
#include "Lucy/Search/IndexSearcher.h"
#include "Lucy/Document/HitDoc.h"

#define DECREF       cfish_Obj_decref
#define LucyIndexSearcher lucy_IndexSearcher
#define LucyIxSearcherNew lucy_IxSearcher_new
#define LucyIxSearcherHits LUCY_IxSearcher_Hits
#define LucyHitDoc lucy_HitDoc
#define LucyHitsTotal LUCY_Hits_Total_Hits
#define LucyHitsNext LUCY_Hits_Next
#define LucyHitDocExtract LUCY_HitDoc_Extract

#include "Lucy/Search/Query.h"
#define LucyQuery lucy_Query

#include "Lucy/Search/QueryParser.h"
#define LucyQueryParser lucy_QueryParser
#define LucyQParserNew lucy_QParser_new
#define LucyQParserParse LUCY_QParser_Parse

#include "Lucy/Plan/Schema.h"
#define LucySchema lucy_Schema

#include "Lucy/Analysis/EasyAnalyzer.h"
#define LucyEasyAnalyzerNew lucy_EasyAnalyzer_new

#include "Lucy/Plan/Schema.h"
#define LucySchemaAllFields LUCY_Schema_All_Fields
*/
import "C"

type Search struct {
	Location     string
	lucySearcher *C.LucyIndexSearcher
}

type Query struct {
	QueryString string
	lucySchema  *C.LucySchema
	lucyQuery   *C.LucyQuery
}

type IndexReader func(*Query, string, uint, uint) (uint, []string)

// This will need to be a bit more generic,
// but for testing this will work fine.
func (search *Search) GetSearcher() IndexReader {
	idxLocation := cb_newf(search.Location)
	search.lucySearcher = C.LucyIxSearcherNew(idxLocation)
	C.DECREF(idxLocation)
	return func(query *Query, field string, offset, limit uint) (uint, []string) {
		getField := cb_newf(field)
		hits := C.LucyIxSearcherHits(search.lucySearcher, query.lucyQuery, C.uint32_t(offset), C.uint32_t(limit), nil)
		totalNumHits := uint(C.LucyHitsTotal(hits))
		requestedNumHits := minUInt(limit, totalNumHits)
		results := make([]string, requestedNumHits)
		var hit *C.LucyHitDoc
		for i := uint(0); i < requestedNumHits; i++ {
			hit = C.LucyHitsNext(hits)
			if hit == nil {
				break
			}
			value_cb := C.LucyHitDocExtract(hit, getField, nil) // do i need to free this
			value := cb_ptr2char(value_cb)                      // do i need to free this
			results[i] = C.GoString(value)
			C.DECREF(hit)
		}
		C.DECREF(getField)
		C.DECREF(hits)
		return totalNumHits, results
	}
}

func (search *Search) Close() {
	C.DECREF(search.lucySearcher)
}

func NewQuery(schema *Schema, queryString string) *Query {
	language := cb_newf("en")
	defer C.DECREF(language)
	analyzer := C.LucyEasyAnalyzerNew(language)
	defer C.DECREF(analyzer)
	qp := C.LucyQParserNew(schema.lucySchema, analyzer, cb_newf("AND"), C.LucySchemaAllFields(schema.lucySchema))
	defer C.DECREF(qp)
	return &Query{
		QueryString: queryString,
		lucySchema:  schema.lucySchema,
		lucyQuery:   C.LucyQParserParse(qp, cb_new_from_utf8(queryString)),
	}
}

func (query *Query) Close() {
	C.DECREF(query.lucyQuery)
}
