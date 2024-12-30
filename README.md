# Ordered Map
[![Go Reference](https://pkg.go.dev/badge/github.com/MarkRosemaker/ordmap.svg)](https://pkg.go.dev/github.com/MarkRosemaker/ordmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarkRosemaker/ordmap)](https://goreportcard.com/report/github.com/MarkRosemaker/ordmap)
![Code Coverage](https://img.shields.io/badge/coverage-98.3%25-green)
[![License: **MIT**](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

`ordmap` is a Go package that provides a generic ordered map implementation, primarily designed for JSON v2 marshalling and unmarshalling using the [go-json-experiment](https://github.com/go-json-experiment/json) library.

An ordered map maintains the order of keys based on insertion, allowing you to iterate over the map in the order in which entries were added. This can be particularly useful for applications where the order of elements is important, such as in JSON serialization or when maintaining the sequence of operations.

## Features

- **Seamless JSON v2 Integration:** Directly integrates with the [JSON v2](https://github.com/go-json-experiment/json) library for efficient and order-preserving marshalling and unmarshalling.
- **Custom Ordered Maps:** Provides robust helper functions to easily define your own custom ordered maps with minimal boilerplate code.
- **Pre-Defined Ordered Map Alias:** Simplifies usage by offering a pre-defined ordered map type that can be conveniently aliased for specific key and value types.
- **Efficient Ordered Operations:** Ensures efficient insertion, retrieval, and iteration while maintaining the order of elements, making it ideal for use cases where order matters.
- **Ordered Iteration:** Leverages the `ByIndex` method to iterate over the map in an ordered manner based on the insertion sequence.

## Installation

To install the library, use the following command:

```shell
go get github.com/MarkRosemaker/ordmap
```

## Usage

### Custom Ordered Map

To create your own custom ordered map, you can utilize helper functions to define its methods:

```go
package main

import (
	"iter"

	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
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

// MarshalJSONV2 marshals the key-value pairs in order.
func (om *MyOrderedMap) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(om, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (om *MyOrderedMap) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(om, dec, opts, setIndex)
}

// Set sets a value in the map, assigning it the highest index + 1.
func (om *MyOrderedMap) Set(key string, v *ValueWithIndex) {
	ordmap.Set(om, key, v, getIndex, setIndex)
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

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/MarkRosemaker/ordmap).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
