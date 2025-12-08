package ordmap_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/MarkRosemaker/ordmap"
)

func (om UserDefinedOrderedMap) ByIndex() iter.Seq2[string, *ValueWithIndex] {
	return ordmap.ByIndex(om, getIndex)
}

func TestByIndex(t *testing.T) {
	t.Parallel()

	t.Run("user defined ordered map", func(t *testing.T) {
		om := UserDefinedOrderedMap{
			"foo": &ValueWithIndex{idx: 1},
			"bar": &ValueWithIndex{idx: 2},
			"baz": &ValueWithIndex{idx: 3},
		}

		testByIndex(t, om, func(s string) { om[s] = &ValueWithIndex{} })
	})

	t.Run("ordered map", func(t *testing.T) {
		var om OrderedMap
		om.Set("foo", Value{})
		om.Set("bar", Value{})
		om.Set("baz", Value{})

		testByIndex[Value](t, om, func(s string) {
			om[s] = ordmap.Value[Value]{}
		})
	})

	t.Run("ordered map with pointer value", func(t *testing.T) {
		var om OrderedMapPointer
		om.Set("foo", &Value{})
		om.Set("bar", &Value{})
		om.Set("baz", &Value{})

		testByIndex[*Value](t, om, func(k string) {
			om[k] = ordmap.Value[*Value]{}
		})
	})
}

func testByIndex[V, Val any, M interface {
	~map[string]Val
	ordmap.ByIndexer[string, V]
}](t *testing.T, om M, set func(string),
) {
	t.Helper()

	for range om.ByIndex() {
		break // i.e. yield returns true
	}

	// test that the keys are sorted by index
	indexedKeys := []string{"foo", "bar", "baz"}

	i := 0
	for k := range om.ByIndex() {
		if indexedKeys[i] != k {
			t.Fatalf("got: %v, want: %v", k, indexedKeys[i])
		}
		i++
	}

	if i != len(indexedKeys) {
		t.Fatalf("got: %d, want: %d", i, len(indexedKeys))
	}

	randomKeys := []string{"qux", "moo", "one", "two", "three"}
	for _, key := range randomKeys {
		set(key)
	}

	// test that additional keys are included but come after the sorted keys
	wantSize := len(indexedKeys) + len(randomKeys)

	i = 0
	for k := range om.ByIndex() {
		if i < len(indexedKeys) {
			if indexedKeys[i] != k {
				t.Fatalf("got: %v, want: %v", k, indexedKeys[i])
			}
		} else {
			if slices.Contains(randomKeys, k) {
				// remove k from want
				randomKeys = slices.DeleteFunc(randomKeys, func(w string) bool { return w == k })
			} else {
				t.Fatalf("unexpected key: %v", k)
			}
		}
		i++
	}

	if len(randomKeys) > 0 {
		t.Fatalf("did not find all random keys, still got: %v", randomKeys)
	}

	if i != wantSize {
		t.Fatalf("got: %d, want: %d", i, wantSize)
	}
}
