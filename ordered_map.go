package ordmap

import (
	"fmt"
	"iter"
	"sort"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"golang.org/x/exp/maps"
)

// OrderedMap is a map that can be ordered.
type OrderedMap[K comparable, V any] interface {
	// ByIndex returns a sequence of key-value pairs sorted by index.
	ByIndex() iter.Seq2[K, V]
}

// JSONV2OrderedMap is an ordered map that can be marshalled and unmarshalled
// using the [JSON v2 package].
//
// [JSON v2 package]: https://pkg.go.dev/github.com/go-json-experiment/json
type JSONV2OrderedMap[K comparable, V any] interface {
	OrderedMap[K, V]
	// OrderedMap must implement json.MarshalerV2 using MarshalOrderedMap.
	json.MarshalerV2
	// OrderedMap must implement json.UnmarshalerV2 using UnmarshalOrderedMap.
	json.UnmarshalerV2
}

// ByIndex is a helper function for an ordered map to implement ByIndex() and fulfill the OrderedMap interface.
func ByIndex[M ~map[K]V, K comparable, V any](m M, getIndex func(V) int) iter.Seq2[K, V] {
	// get the keys and sort them by index
	keys := maps.Keys(m)
	sort.Slice(keys, func(i, j int) bool {
		idxI := getIndex(m[keys[i]])
		idxJ := getIndex(m[keys[j]])
		return idxI != 0 && // if i is not initialized, it should be at the end
			(idxJ == 0 || // if j is not initialized, it should be at the end
				idxI < idxJ) // otherwise, sort by index
	})

	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// UnmarshalJSONV2 is a helper function to implement UnmarshalJSON on an ordered map.
// The setIndex function is called for each value in the map, so that its index is set accordingly.
func UnmarshalJSONV2[M ~map[K]*R, K comparable, R any](
	m *M, dec *jsontext.Decoder, opts json.Options,
	setIndex func(*R, int),
) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	if tkn.Kind() != '{' {
		return fmt.Errorf("expected {, got %s", tkn.Kind())
	}

	// create the map
	*m = M{}

	i := 1 // start at 1 to avoid confusion with zero values

	for {
		// check if we reached the end of the object
		if dec.PeekKind() == '}' {
			_, err := dec.ReadToken() // consume '}', should not fail
			return err
		}

		var key K
		if err := json.UnmarshalDecode(dec, &key, opts); err != nil {
			return err
		}

		var v R
		if err := json.UnmarshalDecode(dec, &v, opts); err != nil {
			return &errpath.ErrKey{Key: fmt.Sprint(key), Err: err}
		}

		// set the index
		setIndex(&v, i)
		i++

		// set the variable in the map
		(*m)[K(key)] = &v
	}
}

// MarshalJSONV2 is a helper function to implement MarshalJSON on an ordered map.
func MarshalJSONV2[M OrderedMap[K, V], K comparable, V any](
	m M, enc *jsontext.Encoder, opts json.Options,
) error {
	if err := enc.WriteToken(jsontext.ObjectStart); err != nil {
		return err // should never fail
	}

	for k, v := range m.ByIndex() {
		if err := json.MarshalEncode(enc, k, opts); err != nil {
			return err
		}

		if err := json.MarshalEncode(enc, v, opts); err != nil {
			return &errpath.ErrKey{Key: fmt.Sprint(k), Err: err}
		}
	}

	return enc.WriteToken(jsontext.ObjectEnd)
}
