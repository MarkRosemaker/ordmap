package ordmap_test

import "github.com/MarkRosemaker/ordmap"

type (
	// a user-defined ordered map
	UserDefinedOrderedMap map[string]*ValueWithIndex
	// an ordered map with a non-pointer value
	OrderedMap = ordmap.OrderedMap[string, Value]
	// an ordered map with a pointer value
	OrderedMapPointer = ordmap.OrderedMap[string, *Value]
)
