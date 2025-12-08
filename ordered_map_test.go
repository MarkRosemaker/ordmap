package ordmap_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/MarkRosemaker/ordmap"
)

var (
	_ json.MarshalerTo     = (*UserDefinedOrderedMap)(nil)
	_ json.UnmarshalerFrom = (*UserDefinedOrderedMap)(nil)
)

type (
	// a user-defined ordered map
	UserDefinedOrderedMap map[string]*ValueWithIndex
	// an ordered map with a non-pointer value
	OrderedMap = ordmap.OrderedMap[string, Value]
	// an ordered map with a pointer value
	OrderedMapPointer = ordmap.OrderedMap[string, *Value]
)

func (om *UserDefinedOrderedMap) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(om, enc)
}

func (om *UserDefinedOrderedMap) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(om, dec,
		func(v *ValueWithIndex, i int) *ValueWithIndex { v.idx = i; return v })
}
