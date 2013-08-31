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

type Document map[string]string
type Documents []Document

func NewDocument() Document {
	return make(Document)
}

func (doc Document) Add(field, value string) {
	doc[field] = value
}

func (doc Document) GetFields() []string {
	docLen := len(doc)
	fields := make([]string, docLen, docLen)
	i := 0
	for key, _ := range doc {
		fields[i] = key
		i++
	}
	return fields
}
