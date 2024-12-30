package ordmap

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Value is a value with an index.
type Value[V any] struct {
	V   V
	idx int
}

// UnmarshalJSONV2 unmarshals a value by just decoding the value.
// The index is set by the caller.
func (cs *Value[_]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return json.UnmarshalDecode(dec, &cs.V, opts)
}

// MarshalJSONV2 marshals a value by encoding just the value and ignoring the index.
func (v Value[_]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return json.MarshalEncode(enc, v.V, opts)
}

// getIndex returns the index of a value.
func getIndex[V any](v Value[V]) int { return v.idx }

// setIndex sets the index of a value.
func setIndex[V any](v *Value[V], i int) { v.idx = i }
