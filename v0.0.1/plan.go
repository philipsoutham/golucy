package golucy

import "runtime"

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
#define CFishObj cfish_Obj
#define DECREF cfish_Obj_decref

#include "Clownfish/CharBuf.h"
#define CFishCharBuf cfish_CharBuf

#include "Clownfish/VArray.h"
#define CFishVArray cfish_VArray
#define CFishVArrayNew  cfish_VA_new
#define CFishVArrayPush CFISH_VA_Push

#include "Lucy/Plan/Schema.h"
#define LucySchema lucy_Schema
#define LucySchemaNew lucy_Schema_new
#define LucySchemaSpecField LUCY_Schema_Spec_Field

#include "Lucy/Analysis/Analyzer.h"
#define LucyAnalyzer lucy_Analyzer

#include "Lucy/Analysis/StandardTokenizer.h"
#define LucyStandardTokenizerNew lucy_StandardTokenizer_new
#define LucyStandardTokenizer lucy_StandardTokenizer

#include "Lucy/Analysis/Normalizer.h"
#define LucyNormalizerNew lucy_Normalizer_new
#define LucyNormalizer lucy_Normalizer

#include "Lucy/Analysis/EasyAnalyzer.h"
#define LucyEasyAnalyzerNew lucy_EasyAnalyzer_new
#define LucyEasyAnalyzer lucy_EasyAnalyzer

#include "Lucy/Analysis/PolyAnalyzer.h"
#define LucyPolyAnalyzerNew lucy_PolyAnalyzer_new
#define LucyPolyAnalyzer lucy_PolyAnalyzer

#include "Lucy/Plan/FullTextType.h"
#define LucyFullTextTypeNew lucy_FullTextType_new
#define LucyFullTextTypeInitOptions lucy_FullTextType_init2

#include "Lucy/Plan/StringType.h"
#define LucyStringTypeNew lucy_StringType_new
#define LucyStringTypeInitOptions lucy_StringType_init2

#include "Lucy/Search/Query.h"
#define LucyQuery lucy_Query
#define LucyMakeCompiler LUCY_Query_Make_Compiler

*/
import "C"

type indexType uint8

const (
	FullTextType indexType = iota
	StringType
)

//type Analyzer struct {
//	Language     string
//	lucyAnalyzer *C.LucyEasyAnalyzer
//}

//type Analyzer struct {
//	Language     string
//	lucyAnalyzer *C.LucyPolyAnalyzer
//}

type Analyzer struct {
	Language     string
	lucyAnalyzer *C.LucyAnalyzer
}

type IndexOptions struct {
	Analyzer      *Analyzer
	Boost         float32
	Indexed       bool
	Stored        bool
	Sortable      bool
	Highlightable bool
}

type Field struct {
	Name         string
	IndexType    indexType
	IndexOptions *IndexOptions
}

type Schema struct {
	Fields     []*Field
	lucySchema *C.LucySchema
}

func NewIdField(name string) *Field {
	return &Field{
		Name:      name,
		IndexType: StringType,
		IndexOptions: &IndexOptions{
			Boost:         0.0,
			Indexed:       false,
			Stored:        true,
			Sortable:      false,
			Highlightable: false,
		},
	}
}

func NewFTField(name, language string, stemTerms bool) *Field {
	return &Field{
		Name:      name,
		IndexType: FullTextType,
		IndexOptions: &IndexOptions{
			Analyzer:      NewAnalyzer(language, stemTerms /* stemming */),
			Boost:         1.0,
			Indexed:       true,
			Stored:        true,
			Sortable:      false,
			Highlightable: true,
		},
	}
}

func NewSchema() *Schema {
	schema := &Schema{lucySchema: C.LucySchemaNew()}
	runtime.SetFinalizer(schema, freeSchema)
	return schema
}

func (schema *Schema) AddField(field *Field) {
	schema.Fields = append(schema.Fields, field)
	var specType *C.CFishObj
	defer C.DECREF(specType)
	name := cb_newf(field.Name)
	defer C.DECREF(name)

	switch field.IndexType {
	case FullTextType:
		specType = fullTextSpecType(field)
	case StringType:
		specType = stringSpecType(field)
	default:
		panic("Specified IndexType not supported yet")
	}
	C.LucySchemaSpecField(schema.lucySchema, name, specType)
}

func (schema *Schema) AddFields(fields ...(*Field)) {
	for _, field := range fields {
		schema.AddField(field)
	}
}

func (schema *Schema) Close() {
	if schema.lucySchema != nil {
		C.DECREF(schema.lucySchema)
		schema.lucySchema = nil
	}
}

func NewIndexOptions(language string, boost float32, indexed, stored, sortable, highlightable bool, stemTerms bool) *IndexOptions {
	return &IndexOptions{
		Analyzer:      NewAnalyzer(language, stemTerms /* stemming */),
		Boost:         boost,
		Indexed:       indexed,
		Stored:        stored,
		Sortable:      sortable,
		Highlightable: highlightable,
	}
}

// keeping this around for posterity
//func NewAnalyzer(language string) *Analyzer {
//	lang := cb_newf(language)
//	defer C.DECREF(lang)
//	analyzer := &Analyzer{Language: language, lucyAnalyzer: C.LucyEasyAnalyzerNew(lang)}
//	fmt.Printf("%+v\n", analyzer)
//	runtime.SetFinalizer(analyzer, freeAnalyzer)
//	return analyzer
//}

func NewAnalyzer(language string, stemTerms bool) *Analyzer {
	lang := cb_newf(language)
	defer C.DECREF(lang)

	// non-stemming analyzer still does case-folding (normalizing) and tokenizing
	var analyzer *Analyzer
	if stemTerms {
		// see https://lucy.apache.org/docs/test/Lucy/Docs/Tutorial/Analysis.html
		analyzer = &Analyzer{Language: language, lucyAnalyzer: C.LucyEasyAnalyzerNew(lang)}
	} else {
		tokenizer := C.LucyStandardTokenizerNew()
		normalizer := C.LucyNormalizerNew(nil, (C.bool)(true), (C.bool)(false))
		analyzers := C.CFishVArrayNew((C.uint32_t)(2))

		//defer C.DECREF(tokenizer) get a segfault if i do this..
		//defer C.DECREF(normalizer) get a segfault if i do this..
		defer C.DECREF(analyzers) // this works, however

		// have to push the tokenizer before the normalizer - otherwise
		// bad bad bad bad bad bad things will happen.
		C.CFishVArrayPush(analyzers, tokenizer)
		C.CFishVArrayPush(analyzers, normalizer)
		analyzer = &Analyzer{Language: language, lucyAnalyzer: C.LucyPolyAnalyzerNew(lang, analyzers)}
	}
	runtime.SetFinalizer(analyzer, freeAnalyzer)
	return analyzer
}

func (analyzer *Analyzer) Close() {
	if analyzer.lucyAnalyzer != nil {
		C.DECREF(analyzer.lucyAnalyzer)
		analyzer.lucyAnalyzer = nil
	}
}

func stringSpecType(field *Field) *C.CFishCharBuf {
	return C.LucyStringTypeInitOptions(
		C.LucyStringTypeNew(),
		(C.float)(field.IndexOptions.Boost),
		(C.bool)(field.IndexOptions.Indexed),
		(C.bool)(field.IndexOptions.Stored),
		(C.bool)(field.IndexOptions.Sortable),
	)
}

func fullTextSpecType(field *Field) *C.CFishCharBuf {
	// Two ways to skin a cat
	//
	// specType := C.LucyFullTextTypeNew(field.IndexOptions.Analyzer.lucyAnalyzer)
	// specType = C.LucyFullTextTypeInitOptions(specType,
	// 	field.IndexOptions.Analyzer.lucyAnalyzer,
	// 	(C.float)(field.IndexOptions.Boost),
	// 	(C.bool)(field.IndexOptions.Indexed),
	// 	(C.bool)(field.IndexOptions.Stored),
	// 	(C.bool)(field.IndexOptions.Sortable),
	// 	(C.bool)(field.IndexOptions.Highlightable),
	// )
	// return specType
	//
	// and another
	return C.LucyFullTextTypeInitOptions(
		C.LucyFullTextTypeNew(field.IndexOptions.Analyzer.lucyAnalyzer),
		field.IndexOptions.Analyzer.lucyAnalyzer,
		(C.float)(field.IndexOptions.Boost),
		(C.bool)(field.IndexOptions.Indexed),
		(C.bool)(field.IndexOptions.Stored),
		(C.bool)(field.IndexOptions.Sortable),
		(C.bool)(field.IndexOptions.Highlightable),
	)
}

func freeSchema(schema *Schema) {
	schema.Close()
}

func freeAnalyzer(analyzer *Analyzer) {
	analyzer.Close()
}
