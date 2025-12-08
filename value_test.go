package ordmap_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/MarkRosemaker/ordmap"
)

type ValueWithIndex struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`

	idx int
}

func getIndex(v *ValueWithIndex) int                    { return v.idx }
func setIndex(v *ValueWithIndex, i int) *ValueWithIndex { v.idx = i; return v }

type Value struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestValue_Arshal(t *testing.T) {
	const want = `{"foo":"foo","bar":1}`

	var v ordmap.Value[Value]
	if err := json.Unmarshal([]byte(want), &v); err != nil {
		t.Fatal(err)
	}

	got, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != want {
		t.Fatalf("got: %v, want: %v", string(got), want)
	}
}
