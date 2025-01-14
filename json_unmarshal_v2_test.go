package ordmap_test

import (
	"reflect"
	"testing"

	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func (om *UserDefinedOrderedMap) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(om, dec, opts, setIndex)
}

func TestUnmarshalJSONFrom(t *testing.T) {
	t.Parallel()

	const want = `{"foo":{"foo":"foo","bar":1},"bar":{"foo":"foo","bar":1},"baz":{"foo":"foo","bar":1},"qux":{"foo":"","bar":0},"moo":{"foo":"","bar":0},"one":{"foo":"","bar":0},"two":{"foo":"","bar":0},"three":{"foo":"","bar":0}}`

	for _, tc := range []struct {
		name string
		om   any
	}{
		{"user defined ordered map", UserDefinedOrderedMap{}},
		{"ordered map", OrderedMap{}},
		{"ordered map with pointer value", OrderedMapPointer{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(want), &tc.om); err != nil {
				t.Fatal(err)
			}

			got, err := json.Marshal(tc.om)
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != want {
				t.Fatalf("got: %v, want: %v", string(got), want)
			}
		})
	}
}

func TestUnmarshal_Errors(t *testing.T) {
	t.Parallel()

	for _, testType := range []struct {
		name string
		tp   reflect.Type
		val  reflect.Type
	}{
		{
			"user defined ordered map",
			reflect.TypeFor[UserDefinedOrderedMap](), reflect.TypeFor[ValueWithIndex](),
		},
		{
			"ordered map",
			reflect.TypeFor[OrderedMap](), reflect.TypeFor[Value](),
		},
		{
			"ordered map with pointer value",
			reflect.TypeFor[OrderedMapPointer](), reflect.TypeFor[Value](),
		},
	} {
		t.Run(testType.name, func(t *testing.T) {
			t.Run("invalid key", func(t *testing.T) {
				err := json.Unmarshal([]byte(`{1}`), reflect.New(testType.tp).Interface())
				synErr := errAs[jsontext.SyntacticError](t, err)
				if synErr.JSONPointer != "" {
					t.Fatalf("got: %v, want: %v", synErr.JSONPointer, "")
				}

				if synErr.Err == nil {
					t.Fatal("expected error")
				} else if want := `object member name must be a string`; synErr.Err.Error() != want {
					t.Fatalf("got: %q, want: %v", synErr.Err, want)
				}
			})

			for _, tc := range []struct {
				name string
				data string
				err  string
			}{
				{"empty", ``, `EOF`},
				{"string instead of object", `""`, `expected {, got string`},
				{
					"missing string for object name", `{"foo":{"foo":"foo","bar":1`,
					`["foo"]: jsontext: unexpected EOF within "/foo" after offset 27`,
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					err := json.Unmarshal([]byte(tc.data), reflect.New(testType.tp).Interface())

					semErr := errAs[json.SemanticError](t, err)
					if semErr.GoType != testType.tp {
						t.Fatalf("got: %v, want: %v", semErr.GoType, testType.tp)
					}

					if semErr.Err == nil {
						t.Fatal("expected error")
					} else if semErr.Err.Error() != tc.err {
						t.Fatalf("got: %q, want: %q", err, tc.err)
					}
				})
			}

			t.Run("missing string for object name", func(t *testing.T) {
				err := json.Unmarshal([]byte(`{"foo":1}`), reflect.New(testType.tp).Interface())

				semErr := errAs[json.SemanticError](t, err)
				if semErr.GoType != testType.tp {
					t.Fatalf("got: %v, want: %v", semErr.GoType, testType.tp)
				}

				semErr = errAs[json.SemanticError](t, semErr.Err)
				if semErr.GoType != testType.val {
					t.Fatalf("got: %v, want: %v", semErr.GoType, testType.val)
				}

				if semErr.Err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}
