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

#include "Lucy/Analysis/EasyAnalyzer.h"
#include "Lucy/Plan/FullTextType.h"
#include "Lucy/Plan/BlobType.h"
#include "Lucy/Plan/StringType.h"
#include "Lucy/Plan/Schema.h"


#define DECREF       cfish_Obj_decref
#define CFishObj cfish_Obj
#define CFishCharBuf cfish_CharBuf

#define LucySchema lucy_Schema
#define LucySchemaNew lucy_Schema_new
#define LucyEasyAnalyzerNew lucy_EasyAnalyzer_new
#define LucyFullTextTypeNew lucy_FullTextType_new
#define LucyFullTextTypeInitOptions lucy_FullTextType_init2
#define LucyBlobTypeNew lucy_BlobType_new
#define LucyStringTypeNew lucy_StringType_new
#define LucyStringTypeInitOptions lucy_StringType_init2
#define LucySchemaSpecField LUCY_Schema_Spec_Field
*/
import "C"

type PlanItemOptions struct {
	Stored        bool
	Indexed       bool
	Sortable      bool
	Highlightable bool
	Boost         float32
	Language      string // used only for FullTextType
}

type PlanItem struct {
	Field   string
	Type    IndexType
	Options *PlanItemOptions
}

type Schema struct {
	PlanItems  []*PlanItem
	lucySchema *C.LucySchema
}

func (schema *Schema) AddPlanItem(item *PlanItem) {
	schema.PlanItems = append(schema.PlanItems, item)
}

func (schema *Schema) Commit() {
	schema.createLucySchema()
}

func (schema *Schema) createLucySchema() {
	lucySchema := C.LucySchemaNew()
	for _, item := range schema.PlanItems {
		var specType *C.CFishObj
		if item.Type == FullTextType {
			var language *C.CFishCharBuf
			if item.Options != nil && item.Options.Language != "" {
				language = cb_newf(item.Options.Language)
			} else {
				language = cb_newf("en")
			}
			analyzer := C.LucyEasyAnalyzerNew(language)
			specType = C.LucyFullTextTypeNew(analyzer)
			// TODO: come up with a better way to handle options.
			// This isn't very friendly.
			if item.Options != nil {
				specType = C.LucyFullTextTypeInitOptions(specType, analyzer,
					(C.float)(item.Options.Boost),
					(C.bool)(item.Options.Indexed),
					(C.bool)(item.Options.Stored),
					(C.bool)(item.Options.Sortable),
					(C.bool)(item.Options.Highlightable),
				)
			}
			C.DECREF(language)
			C.DECREF(analyzer)
		} else if item.Type == StringType {
			specType = C.LucyStringTypeNew()
			if item.Options != nil {
				specType = C.LucyStringTypeInitOptions(specType,
					(C.float)(item.Options.Boost),
					(C.bool)(item.Options.Indexed),
					(C.bool)(item.Options.Stored),
					(C.bool)(item.Options.Sortable),
				)
			}
		} else if item.Type == BlobType {
			isStored := (C.bool)(false)
			if item.Options != nil && item.Options.Stored {
				isStored = (C.bool)(true)
			}
			specType = C.LucyBlobTypeNew(isStored)
			// need to send []cfish_byte castable value
			panic("BlobType not supported yet")
		} else {
			panic("Type not supported yet")
		}
		fieldName := cb_newf(item.Field)
		C.LucySchemaSpecField(lucySchema, fieldName, specType)
		C.DECREF(fieldName)
		C.DECREF(specType)
	}
	schema.lucySchema = lucySchema
}

func (schema *Schema) Close() {
	C.DECREF(schema.lucySchema)
}
