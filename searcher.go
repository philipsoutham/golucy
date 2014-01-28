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
#define ObjToString CFISH_Obj_To_String

#include "Clownfish/Hash.h"
#define CfishHashKeys CFISH_Hash_Keys //returns cfish_Varray*

#include "Lucy/Search/IndexSearcher.h"
#define LucyIndexSearcher lucy_IndexSearcher
#define LucyIxSearcherNew lucy_IxSearcher_new
#define LucyIxSearcherHits LUCY_IxSearcher_Hits
#define LucyIxSearcherGetSchema LUCY_IxSearcher_Get_Schema
#define LucyIxSearchFetchDocVec LUCY_IxSearcher_Fetch_Doc_Vec

#include "Lucy/Analysis/EasyAnalyzer.h"
#define LucyEasyAnalyzerNew lucy_EasyAnalyzer_new
#define LucyEasyAnalyzerTransform LUCY_EasyAnalyzer_Transform // returns lucy_Inversion*
#define LucyEasyAnalyzerTransformText LUCY_EasyAnalyzer_Transform_Text // returns lucy_Inversion*

#include "Lucy/Analysis/Token.h"
#define LucyInversionNext LUCY_Inversion_Next //returns lucy_Token*
#define LucyInversionNextCluster LUCY_Inversion_Next_Cluster //returns lucy_Token**



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
#define LucyHitDocGetDocId  LUCY_HitDoc_Get_Doc_ID
#define LucyHitdocDump LUCY_HitDoc_Dump //returns cfish_Hash*
#define LucyHitDocGetScore LUCY_HitDoc_Get_Score

#include "Lucy/Search/Hits.h"
#define LucyHitsTotal LUCY_Hits_Total_Hits
#define LucyHitsNext LUCY_Hits_Next

#include "Lucy/Search/Query.h"
#define LucyQuery lucy_Query
#define LucyQueryMakeCompiler LUCY_Query_Make_Compiler

#include "Lucy/Search/Compiler.h"
#define LucyCompilerHighlightSpans LUCY_Compiler_Highlight_Spans

#include "Clownfish/VArray.h"
#define VaGetSize CFISH_VA_Get_Size
#define VaFetch CFISH_VA_Fetch

#include "Lucy/Search/Span.h"
#define LucySpanGetOffset LUCY_Span_Get_Offset
#define LucySpanGetLength LUCY_Span_Get_Length
*/
import "C"

import "strings"

type Query struct {
	QueryStr  string
	lucyQuery *C.LucyQuery
}

type SearchResult struct {
	Id           string
	Text         string
	Score        float32
	MatchedTerms []string
}

type IndexReader struct {
	Index        *Index
	lucySearcher *C.LucyIndexSearcher
}

func (index *Index) NewIndexReader() *IndexReader {
	ixLocation := cb_newf(index.Path)
	defer C.DECREF(ixLocation)
	ixReader := &IndexReader{
		Index:        index,
		lucySearcher: C.LucyIxSearcherNew(ixLocation),
	}
	runtime.SetFinalizer(ixReader, freeIndexReader)
	return ixReader
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
	qs := cb_new_from_utf8(queryStr)
	defer C.DECREF(qs)
	query := &Query{
		QueryStr:  queryStr,
		lucyQuery: C.LucyQParserParse(qp, qs),
	}
	runtime.SetFinalizer(query, freeQuery)
	return query
}

func (ixReader *IndexReader) Search(query *Query, offset, limit uint, idField string, contentField string, includeMatchedTerms bool) (uint, []*SearchResult) {
	// Should probably have some sort
	// of `Results` object/iterator so that we don't have to specify
	// offset/limit and where I can attach matched terms to the result.
	lIdField, lContentField := cb_newf(idField), cb_newf(contentField) // total hack, need to return more than one field
	defer C.DECREF(lIdField)
	defer C.DECREF(lContentField)
	hits := C.LucyIxSearcherHits(ixReader.lucySearcher, query.lucyQuery, C.uint32_t(offset), C.uint32_t(limit), nil)
	defer C.DECREF(hits)
	totalNumHits := uint(C.LucyHitsTotal(hits))
	num2Return := minUInt(limit, totalNumHits)
	results := make([]*SearchResult, num2Return)
	var hit *C.LucyHitDoc
	compiler := C.LucyQueryMakeCompiler(query.lucyQuery, ixReader.lucySearcher, 1.0, false)
	defer C.DECREF(compiler)

	matchedTerms := func(docId C.int32_t, result *SearchResult) {
		docVec := C.LucyIxSearchFetchDocVec(ixReader.lucySearcher, docId)
		defer C.DECREF(docVec)
		spans := C.LucyCompilerHighlightSpans(compiler, ixReader.lucySearcher, docVec, lContentField)
		defer C.DECREF(spans)
		spanCnt := C.VaGetSize(spans)
		if spanCnt == 0 {
			// should never get here, but just in case...
			return
		}
		result.MatchedTerms = make([]string, spanCnt)
		var i C.uint32_t
		for i = 0; i < spanCnt; i++ {
			span := C.VaFetch(spans, i)
			offset := C.LucySpanGetOffset(span)
			length := C.LucySpanGetLength(span)
			result.MatchedTerms[i] = string([]rune(result.Text)[offset : offset+length])
		}
		// make terms unique?
		result.MatchedTerms = set(result.MatchedTerms)
	}
	var i uint
	for i = 0; i < num2Return; i++ {
		hit = C.LucyHitsNext(hits)
		if hit == nil {
			break
		}
		docId := C.LucyHitDocGetDocId(hit)
		contentValue := cb_ptr2char(C.LucyHitDocExtract(hit, lContentField)) // do i need to free this
		idValue := cb_ptr2char(C.LucyHitDocExtract(hit, lIdField))           // do i need to free this
		results[i] = &SearchResult{
			Id:    C.GoString(idValue),
			Text:  C.GoString(contentValue),
			Score: float32(C.LucyHitDocGetScore(hit)),
		}
		if includeMatchedTerms {
			matchedTerms(docId, results[i])
		}
		C.DECREF(hit)
	}
	return totalNumHits, results
}

func set(vals []string) []string {
	s := make(map[string]bool)
	for _, val := range vals {
		s[strings.ToLower(val)] = true
	}
	retval := make([]string, len(s))
	i := 0
	for k, _ := range s {
		retval[i] = k
		i++
	}
	return retval
}

func (ixReader *IndexReader) Close() {
	if ixReader.lucySearcher != nil {
		C.DECREF(ixReader.lucySearcher)
		ixReader.lucySearcher = nil
	}
}

func (query *Query) Close() {
	if query.lucyQuery != nil {
		C.DECREF(query.lucyQuery)
		query.lucyQuery = nil
	}
}

func freeIndexReader(ixReader *IndexReader) {
	ixReader.Close()
}

func freeQuery(query *Query) {
	query.Close()
}
