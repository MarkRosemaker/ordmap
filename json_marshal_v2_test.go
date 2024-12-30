package ordmap_test

import (
	"errors"
	"iter"
	"reflect"
	"testing"

	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func (om *UserDefinedOrderedMap) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(om, enc, opts)
}

func TestMarshalJSONV2(t *testing.T) {
	t.Parallel()

	const want = `{"foo":{"foo":"a","bar":6},"bar":{"foo":"b","bar":7},"baz":{"foo":"c","bar":8}}`

	t.Run("user defined ordered map", func(t *testing.T) {
		om := UserDefinedOrderedMap{
			"foo": &ValueWithIndex{Foo: "a", Bar: 6, idx: 1},
			"bar": &ValueWithIndex{Foo: "b", Bar: 7, idx: 2},
		}

		om["baz"] = &ValueWithIndex{Foo: "c", Bar: 8}

		got, err := json.Marshal(om)
		if err != nil {
			t.Fatal(err)
		}

		if want != string(got) {
			t.Fatalf("got: %v, want: %v", string(got), want)
		}

		got, err = json.Marshal(UserDefinedOrderedMap{
			"qux": nil,
		})
		if err != nil {
			t.Fatal(err)
		}

		if want := `{"qux":null}`; want != string(got) {
			t.Fatalf("got: %v, want: %v", string(got), want)
		}
	})

	t.Run("ordered map", func(t *testing.T) {
		var om OrderedMap
		om.Set("foo", Value{Foo: "a", Bar: 6})
		om.Set("bar", Value{Foo: "b", Bar: 7})
		om["baz"] = ordmap.Value[Value]{V: Value{Foo: "c", Bar: 8}}

		got, err := json.Marshal(om)
		if err != nil {
			t.Fatal(err)
		}

		if want != string(got) {
			t.Fatalf("got: %v, want: %v", string(got), want)
		}

		got, err = json.Marshal(OrderedMap{
			"qux": ordmap.Value[Value]{},
		})
		if err != nil {
			t.Fatal(err)
		}

		if want := `{"qux":{"foo":"","bar":0}}`; want != string(got) {
			t.Fatalf("got: %v, want: %v", string(got), want)
		}
	})

	t.Run("ordered map with pointer value", func(t *testing.T) {
		var om OrderedMapPointer
		om.Set("foo", &Value{Foo: "a", Bar: 6})
		om.Set("bar", &Value{Foo: "b", Bar: 7})
		om["baz"] = ordmap.Value[*Value]{V: &Value{Foo: "c", Bar: 8}}

		got, err := json.Marshal(om)
		if err != nil {
			t.Fatal(err)
		}

		if want != string(got) {
			t.Fatalf("got: %v, want: %v", string(got), want)
		}

		got, err = json.Marshal(OrderedMapPointer{
			"qux": ordmap.Value[*Value]{},
		})
		if err != nil {
			t.Fatal(err)
		}

		if want := `{"qux":null}`; want != string(got) {
			t.Fatalf("got: %v, want: %v", string(got), want)
		}
	})
}

// a struct that cannot be marshalled
type impossibleToMarshal struct{ err error }

func (imp impossibleToMarshal) MarshalJSONV2(*jsontext.Encoder, json.Options) error {
	return imp.err
}

// a map which has keys that cannot be marshalled
type cannotMarshalKey map[impossibleToMarshal]*ValueWithIndex

func (om cannotMarshalKey) ByIndex() iter.Seq2[impossibleToMarshal, *ValueWithIndex] {
	return ordmap.ByIndex(om, func(v *ValueWithIndex) int { return v.idx })
}

func (om cannotMarshalKey) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(&om, enc, opts)
}

func (om *cannotMarshalKey) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(om, dec, opts, func(v *ValueWithIndex, i int) { v.idx = i })
}

// a map which has values that cannot be marshalled
type cannotMarshalValue map[string]*impossibleToMarshal

func (om cannotMarshalValue) ByIndex() iter.Seq2[string, *impossibleToMarshal] {
	return ordmap.ByIndex(om, func(v *impossibleToMarshal) int { return 0 })
}

func (om cannotMarshalValue) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(&om, enc, opts)
}

func (om *cannotMarshalValue) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(om, dec, opts, func(v *impossibleToMarshal, i int) {})
}

func TestMarshalJSONV2_Errors(t *testing.T) {
	t.Parallel()

	someErr := errors.New("some error")

	t.Run("marshalling key", func(t *testing.T) {
		_, err := json.Marshal(cannotMarshalKey{
			impossibleToMarshal{err: someErr}: &ValueWithIndex{},
		})

		semErr := errAs[json.SemanticError](t, err)
		if goType := reflect.TypeFor[impossibleToMarshal](); semErr.GoType != goType {
			t.Fatalf("got: %v, want: %v", semErr.GoType, goType)
		} else if semErr.Err != someErr {
			t.Fatalf("got: %v, want: %v", err, someErr)
		}
	})

	t.Run("marshalling value", func(t *testing.T) {
		_, err := json.Marshal(cannotMarshalValue{
			"foo": &impossibleToMarshal{err: someErr},
		})

		semErr := errAs[json.SemanticError](t, err)
		if goType := reflect.TypeFor[cannotMarshalValue](); semErr.GoType != goType {
			t.Fatalf("got: %v, want: %v", semErr.GoType, goType)
		}

		semErr = errAs[json.SemanticError](t, semErr.Err)
		if goType := reflect.TypeFor[impossibleToMarshal](); semErr.GoType != goType {
			t.Fatalf("got: %v, want: %v", semErr.GoType, goType)
		} else if semErr.Err != someErr {
			t.Fatalf("got: %v, want: %v", err, someErr)
		}
	})
}

func errAs[T any, E interface {
	*T
	error
}](t *testing.T, err error,
) E {
	t.Helper()

	var zero T
	target := E(&zero)
	if !errors.As(err, &target) {
		t.Fatalf("want: %T, got: %T", target, err)
	}

	return target
}
