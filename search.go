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

*/
import "C"

type Search struct {
	Location     string
	lucySearcher *C.LucyIndexSearcher
}

type IndexReader func(string, string, uint, uint) (uint, []string)

// This will need to be a bit more generic,
// but for testing this will work fine.
func (search *Search) GetSearcher() IndexReader {
	idxLocation := cb_newf(search.Location)
	search.lucySearcher = C.LucyIxSearcherNew(idxLocation)
	C.DECREF(idxLocation)
	return func(q string, field string, offset, limit uint) (uint, []string) {
		query := cb_new_from_utf8(q)
		getField := cb_newf(field)
		hits := C.LucyIxSearcherHits(search.lucySearcher, query, C.uint32_t(offset), C.uint32_t(limit), nil)
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
		C.DECREF(query)
		C.DECREF(getField)
		C.DECREF(hits)
		return totalNumHits, results
	}
}

func (search *Search) Close() {
	C.DECREF(search.lucySearcher)
}
