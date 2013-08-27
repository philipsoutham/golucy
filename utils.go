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
#include "Clownfish/CharBuf.h"

#define CFishCharBuf cfish_CharBuf
#define DECREF       cfish_Obj_decref


extern CFishCharBuf* CB_newf(const char* pattern) {
    return cfish_CB_newf(pattern);
}

extern char* cfish_cb_ptr2char(const CFishCharBuf * field) {
    return (char*)CFISH_CB_Get_Ptr8(field);
}

*/
import "C"

import (
	"unsafe"
)

func cb_newf(s string) *C.CFishCharBuf {
	cString := C.CString(s)
	defer C.free(unsafe.Pointer(cString))
	return C.CB_newf(cString)
}

func cb_new_from_utf8(s string) *C.CFishCharBuf {
	val := C.CString(s)
	defer C.free(unsafe.Pointer(val))
	vlen := len(s)
	return C.cfish_CB_new_from_utf8(val, (C.size_t)(vlen))
}

func cb_ptr2char(field *C.CFishCharBuf) *C.char {
	return C.cfish_cb_ptr2char(field)
}

func minUInt(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}
