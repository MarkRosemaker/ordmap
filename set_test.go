package ordmap_test

import (
	"testing"

	"github.com/MarkRosemaker/ordmap"
)

func (om *UserDefinedOrderedMap) Set(key string, v *ValueWithIndex) {
	ordmap.Set(om, key, v, func(v *ValueWithIndex) int { return v.idx }, func(v **ValueWithIndex, i int) { (*v).idx = i }) // TODO
}

func TestSet(t *testing.T) {
	om := UserDefinedOrderedMap{
		"foo": &ValueWithIndex{Foo: "a", Bar: 6, idx: 5},
		"bar": &ValueWithIndex{Foo: "b", Bar: 7, idx: 10},
	}

	om.Set("baz", &ValueWithIndex{Foo: "c", Bar: 8})

	if len(om) != 3 {
		t.Fatalf("got: %v, want: 3", len(om))
	}

	i := 0
	for _, v := range om.ByIndex() {
		switch i {
		case 0:
			if v.Foo != "a" || v.Bar != 6 || v.idx != 5 {
				t.Fatalf("got: %#v", v)
			}
		case 1:
			if v.Foo != "b" || v.Bar != 7 || v.idx != 10 {
				t.Fatalf("got: %#v", v)
			}
		case 2:
			if v.Foo != "c" || v.Bar != 8 || v.idx != 11 {
				t.Fatalf("got: %#v", v)
			}
		}

		i++
	}
}
