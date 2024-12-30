package ordmap_test

import (
	"cmp"
	"testing"

	"github.com/MarkRosemaker/ordmap"
)

func (om UserDefinedOrderedMap) Sort() {
	ordmap.Sort(om, setIndex)
}

func (om UserDefinedOrderedMap) SortFunc(less func(string, string) int) {
	ordmap.SortFunc(om, setIndex, less)
}

func less[K cmp.Ordered](a, b K) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}

func TestSort(t *testing.T) {
	t.Parallel()

	t.Run("user defined ordered map", func(t *testing.T) {
		var om UserDefinedOrderedMap

		om.Sort() // no panic

		om = UserDefinedOrderedMap{
			"c": &ValueWithIndex{Foo: "c", idx: 1},
			"a": &ValueWithIndex{Foo: "a", idx: 2},
			"b": &ValueWithIndex{Foo: "b", idx: 3},
		}

		testSort(t, om, om.Sort)

		om = UserDefinedOrderedMap{
			"c": &ValueWithIndex{Foo: "c", idx: 1},
			"a": &ValueWithIndex{Foo: "a", idx: 2},
			"b": &ValueWithIndex{Foo: "b", idx: 3},
		}

		testSort(t, om, func() { om.SortFunc(less) })
	})

	t.Run("ordered map", func(t *testing.T) {
		var om OrderedMap

		om.Sort(less) // no panic

		om.Set("c", Value{Foo: "c"})
		om.Set("a", Value{Foo: "a"})
		om["b"] = ordmap.Value[Value]{V: Value{Foo: "b"}}

		testSort(t, om, func() { om.Sort(less) })
	})

	t.Run("ordered map with pointer value", func(t *testing.T) {
		var om OrderedMapPointer

		om.Sort(less) // no panic

		om.Set("c", &Value{Foo: "c"})
		om.Set("a", &Value{Foo: "a"})
		om["b"] = ordmap.Value[*Value]{V: &Value{Foo: "b"}}

		testSort(t, &om, func() { om.Sort(less) })
	})
}

func testSort[V any](t *testing.T, om ordmap.ByIndexer[string, V], sort func()) {
	t.Helper()

	keys := []string{"c", "a", "b"}

	i := 0
	for k := range om.ByIndex() {
		if k != keys[i] {
			t.Fatalf("got: %v, want: %v", k, keys[i])
		}

		i++
	}

	sort()

	keysSorted := []string{"a", "b", "c"}

	i = 0
	for k := range om.ByIndex() {
		if k != keysSorted[i] {
			t.Fatalf("got: %v, want: %v", k, keysSorted[i])
		}

		i++
	}
}
