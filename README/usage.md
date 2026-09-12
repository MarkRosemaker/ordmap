To install the library, use the following command:

```shell
go get github.com/MarkRosemaker/ordmap
```


### Custom Ordered Map

To create your own custom ordered map, you can utilize helper functions to define its methods:

```go
package main

import (
	"encoding/json/v2"
	"encoding/json/jsontext"
	"iter"

	"github.com/MarkRosemaker/ordmap"
)

type MyOrderedMap map[string]*ValueWithIndex

type ValueWithIndex struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`

	idx int // to order a map of this type
}

func getIndex(v *ValueWithIndex) int                    { return v.idx }
func setIndex(v *ValueWithIndex, i int) *ValueWithIndex { v.idx = i; return v }

// ByIndex returns a sequence of key-value pairs ordered by index.
func (om MyOrderedMap) ByIndex() iter.Seq2[string, *ValueWithIndex] {
	return ordmap.ByIndex(om, getIndex)
}

// Sort sorts the map by key and sets the indices accordingly.
func (om MyOrderedMap) Sort() {
	ordmap.Sort(om, setIndex)
}

// Set sets a value in the map, adding it at the end of the order.
func (om *MyOrderedMap) Set(key string, v *ValueWithIndex) {
	ordmap.Set(om, key, v, getIndex, setIndex)
}

// MarshalJSONTo marshals the key-value pairs in order.
func (om *MyOrderedMap) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONTo(om, enc, opts)
}

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (om *MyOrderedMap) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(om, dec, opts, setIndex)
}
```

If you prefer the map values to be non-pointer types, you can adjust the implementation as follows:

```go
type MyOrderedMap map[string]ValueWithIndex

func getIndex(v ValueWithIndex) int                   { return v.idx }
func setIndex(v ValueWithIndex, i int) ValueWithIndex { v.idx = i; return v }

func (om MyOrderedMap) ByIndex() iter.Seq2[string, ValueWithIndex] {
	return ordmap.ByIndex(om, getIndex)
}

func (om *MyOrderedMap) Set(key string, v ValueWithIndex) {
	ordmap.Set(om, key, v, getIndex, setIndex)
}
```

### Using The Pre-Defined Ordered Map

For simplicity, an ordered map type is already defined for you. You only need to specify the key and value types:

```go
package main

import (
	"github.com/MarkRosemaker/ordmap"
)

type MyOrderedMap = ordmap.OrderedMap[string, *MyValue]

type MyValue struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}
```
