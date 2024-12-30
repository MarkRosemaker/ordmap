package ordmap

import (
	"fmt"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func (om *OrderedMap[K, V]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	// return UnmarshalJSONV2(om, dec, opts, setIndex) TODO

	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	if tkn.Kind() != '{' {
		return fmt.Errorf("expected {, got %s", tkn.Kind())
	}

	// create the map
	*om = OrderedMap[K, V]{}

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

		var v V
		if err := json.UnmarshalDecode(dec, &v, opts); err != nil {
			return &errpath.ErrKey{Key: fmt.Sprint(key), Err: err}
		}

		// set the variable in the map with the correct index
		(*om)[K(key)] = Value[V]{V: v, idx: i}
		i++
	}
}

// UnmarshalJSONV2 is a helper function to make an ordered map fulfil the `json.UnmarshalerV2` interface.
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
