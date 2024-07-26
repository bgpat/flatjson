package flatjson_test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/bgpat/flatjson"
	"github.com/google/go-cmp/cmp"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		input any
		want  flatjson.FlatJSON
		err   error
	}{
		{
			input: map[string]any{"a": 1, "b": map[string]any{"c": 2}},
			want: flatjson.FlatJSON{
				{Path: flatjson.Path{}, Value: map[string]any{"a": 1, "b": map[string]any{"c": 2}}},
				{Path: flatjson.Path{"a"}, Value: 1},
				{Path: flatjson.Path{"b"}, Value: map[string]any{"c": 2}},
				{Path: flatjson.Path{"b", "c"}, Value: 2},
			},
		},
		{
			input: []any{1, 2, map[string]any{"a": 3}},
			want: flatjson.FlatJSON{
				{Path: flatjson.Path{}, Value: []any{1, 2, map[string]any{"a": 3}}},
				{Path: flatjson.Path{0}, Value: 1},
				{Path: flatjson.Path{1}, Value: 2},
				{Path: flatjson.Path{2}, Value: map[string]any{"a": 3}},
				{Path: flatjson.Path{2, "a"}, Value: 3},
			},
		},
		{
			input: make(chan int), // Non-serializable type
			err:   &json.UnsupportedTypeError{Type: reflect.TypeFor[chan int]()},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := flatjson.Flatten(tt.input)
			if err != nil && tt.err == nil || err == nil && tt.err != nil || err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("Flatten() returned an error: %q, want %q", err, tt.err)
			}
			gotJSON, err := json.MarshalIndent(got, "", "  ")
			if err != nil {
				t.Error(err)
			}
			wantJSON, err := json.MarshalIndent(tt.want, "", "  ")
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(gotJSON, wantJSON); diff != "" {
				t.Errorf("Flatten() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		input flatjson.FlatJSON
		path  flatjson.Path
		want  any
		found bool
	}{
		{
			input: flatjson.FlatJSON{
				{Path: flatjson.Path{"a"}, Value: 1},
				{Path: flatjson.Path{"b", "c"}, Value: 2},
			},
			path:  flatjson.Path{"a"},
			want:  1,
			found: true,
		},
		{
			input: flatjson.FlatJSON{
				{Path: flatjson.Path{"a"}, Value: 1},
				{Path: flatjson.Path{"b", "c"}, Value: 2},
			},
			path:  flatjson.Path{"b", "c"},
			want:  2,
			found: true,
		},
		{
			input: flatjson.FlatJSON{
				{Path: flatjson.Path{"a"}, Value: 1},
				{Path: flatjson.Path{"b", "c"}, Value: 2},
			},
			path:  flatjson.Path{"d"},
			want:  nil,
			found: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, found := tt.input.Get(tt.path)
			if found != tt.found {
				msg := "found"
				if !found {
					msg = "not found"
				}
				t.Errorf("Get(%#v) is %s", tt.path, msg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get(%#v) = %#v, want %#v", tt.path, got, tt.want)
			}
		})
	}
}

func ExampleFlatten() {
	obj := map[string]any{
		"a": true,
		"b": map[string]any{
			"c": 2,
		},
		"d": []string{
			"foo",
			"bar",
		},
	}
	flat, err := flatjson.Flatten(obj)
	if err != nil {
		log.Fatal(err)
	}
	for _, pv := range flat {
		path, err := json.Marshal(pv.Path)
		if err != nil {
			log.Fatal(err)
		}
		value, err := json.Marshal(pv.Value)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s = %s\n", path, value)
	}
	// Output:
	// [] = {"a":true,"b":{"c":2},"d":["foo","bar"]}
	// ["a"] = true
	// ["b"] = {"c":2}
	// ["b","c"] = 2
	// ["d"] = ["foo","bar"]
	// ["d",0] = "foo"
	// ["d",1] = "bar"

}

func ExampleFlatJSON_Get() {
	flat := flatjson.FlatJSON{
		{Path: flatjson.Path{"a"}, Value: 1},
		{Path: flatjson.Path{"b", "c"}, Value: 2},
	}
	if value, found := flat.Get(flatjson.Path{"b", "c"}); found {
		fmt.Println(value)
	}
	// Output: 2
}
