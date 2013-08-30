# Overview

Some thoughts on API design and "object" relationships.

```
       Document |
Schema -> Index |-> IndexWriter
             |---> IndexReader -> Searcher
```


## 1. Schema

A `Schema` is comprized of  one to many `Field` types describing it's structure.

```
------|
Field |\
Field | ----> Schema
Field |/
------|
```


### 1.1. Field

A `Field` has a `Name` and  `Type`. The `Name` is a lablel for how the data will be referenced, and the `IndexedType`, in combination with the `IndexOptions`, indicates how it will be treated by the `IndexWriter` and `Searcher`.

```
-------------|
Name         |
IndexType    | ----> Field
IndexOptions |
-------------|

```


### 1.1.1. Name

Just a string value


### 1.1.2. IndexType

An enum...


### 1.1.3. IndexOptions

```
--------------|
Analyzer      |
Boost         |
Indexed       |
Stored        |
Sortable      |
Highlightable |
--------------| ----> IndexOptions
```

#### 1.1.3.1. Analyzer

An Analyzer is a filter which processes text, transforming it from one form into another. For instance, an analyzer might break up a long text into smaller pieces (RegexTokenizer), or it might perform case folding to facilitate case-insensitive search (CaseFolder).


#### 1.1.3.2. Boost

Floating point per-field boost.


#### 1.1.3.3. Indexed

Boolean indicating whether the field should be indexed.


#### 1.1.3.4. Stored

Boolean indicating whether the field should be stored.


#### 1.1.3.5. Highlightable

Boolean indicating whether the field should be highlightable. This also seems to have an impact as to whether or not I will be able to include functionality for extracting matched terms out of matched documents.


## 2. Index

Contains a `Location` and [`Schema`](#1-schema). At a highlevel, describes what is contained in the index files and where they're located.

```
---------|
Location |
Schema   |
---------| ----> Index
```


###  2.1. Location

A string value indicating the path for the directory where the necessary files will be stored.


### 2.2. Schema

Described [here](#1-schema)

<hr/>
<blockquote>
***Note (2013-08-29):***<br/>
*My document writing enthusiasm is starting to dwindle, expect laziness from this point on.*
</blockquote>
<hr/>

## 3. Document

For right now, a `map[string][string]` may be the best way to reflect this data type.
In Go:
```go
type Document map[string]string

func NewDocument() Document {
	return make(Document)
}

func (doc Document) Add(field, value string) {
    doc[field] = value
}
```

## 4. IndexWriter

How documents get added to the index.

```go
func (iw *IndexWriter) AddDoc(doc Document) {
}

func (iw *IndexWriter) AddDocs(docs ...Document) {
}
```
