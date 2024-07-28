package flatjson_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bgpat/flatjson"
)

var (
	json1 = map[string]any{
		"name": "Alice",
		"age":  30,
		"address": map[string]any{
			"city": "Tokyo",
			"zip":  "100-0001",
		},
	}
	json2 = map[string]any{
		"name": "Alice",
		"age":  31,
		"address": map[string]any{
			"city":    "Kyoto",
			"country": "Japan",
		},
	}
)

func TestDiff(t *testing.T) {
	tests := map[string]struct {
		source any
		target any
		want   []flatjson.DiffOperation
		err    error
	}{
		"no changes": {
			source: json1,
			target: json1,
			want:   nil,
		},
		"diff": {
			source: json1,
			target: json2,
			want: []flatjson.DiffOperation{
				{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{}, Value: json2, OldValue: json1},
				{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"address"}, OldValue: map[string]string{"city": "Tokyo", "zip": "100-0001"}, Value: map[string]string{"city": "Kyoto", "country": "Japan"}},
				{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"address", "city"}, Value: "Kyoto", OldValue: "Tokyo"},
				{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"address", "country"}, Value: "Japan"},
				{Type: flatjson.DiffOperationTypeRemove, Path: flatjson.Path{"address", "zip"}},
				{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"age"}, Value: 31, OldValue: 30},
			},
		},
		"error": {
			source: json1,
			target: make(chan int), // Non-serializable type
			err:    &json.UnsupportedTypeError{Type: reflect.TypeFor[chan int]()},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := flatjson.Diff(tt.source, tt.target)
			if err != nil && tt.err == nil || err == nil && tt.err != nil || err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("Diff() returned an error: %q, want %q", err, tt.err)
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
				t.Errorf("Diff() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiffOperation_MarshalJSON(t *testing.T) {
	want := `[
  {
    "path": [
      "add",
      "empty"
    ],
    "type": "add",
    "value": ""
  },
  {
    "path": [
      "add",
      "false"
    ],
    "type": "add",
    "value": false
  },
  {
    "path": [
      "add",
      "null"
    ],
    "type": "add",
    "value": null
  },
  {
    "path": [
      "add",
      "number"
    ],
    "type": "add",
    "value": 1
  },
  {
    "path": [
      "add",
      "string"
    ],
    "type": "add",
    "value": "a"
  },
  {
    "path": [
      "add",
      "true"
    ],
    "type": "add",
    "value": true
  },
  {
    "path": [
      "add",
      "zero"
    ],
    "type": "add",
    "value": 0
  },
  {
    "path": [
      "remove"
    ],
    "type": "remove"
  },
  {
    "old_value": "",
    "path": [
      "replace",
      "empty"
    ],
    "type": "replace",
    "value": "a"
  },
  {
    "old_value": false,
    "path": [
      "replace",
      "false"
    ],
    "type": "replace",
    "value": true
  },
  {
    "old_value": null,
    "path": [
      "replace",
      "null"
    ],
    "type": "replace",
    "value": 1
  },
  {
    "old_value": 1,
    "path": [
      "replace",
      "number"
    ],
    "type": "replace",
    "value": 0
  },
  {
    "old_value": "a",
    "path": [
      "replace",
      "string"
    ],
    "type": "replace",
    "value": ""
  },
  {
    "old_value": true,
    "path": [
      "replace",
      "true"
    ],
    "type": "replace",
    "value": false
  },
  {
    "old_value": 0,
    "path": [
      "replace",
      "zero"
    ],
    "type": "replace",
    "value": 1
  }
]`
	got, err := json.MarshalIndent([]flatjson.DiffOperation{
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "empty"}, Value: ""},
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "false"}, Value: false},
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "null"}, Value: nil},
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "number"}, Value: 1},
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "string"}, Value: "a"},
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "true"}, Value: true},
		{Type: flatjson.DiffOperationTypeAdd, Path: flatjson.Path{"add", "zero"}, Value: 0},
		{Type: flatjson.DiffOperationTypeRemove, Path: flatjson.Path{"remove"}},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "empty"}, Value: "a", OldValue: ""},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "false"}, Value: true, OldValue: false},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "null"}, Value: 1, OldValue: nil},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "number"}, Value: 0, OldValue: 1},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "string"}, Value: "", OldValue: "a"},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "true"}, Value: false, OldValue: true},
		{Type: flatjson.DiffOperationTypeReplace, Path: flatjson.Path{"replace", "zero"}, Value: 1, OldValue: 0},
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(got), want); diff != "" {
		t.Errorf("MarshalJSON() mismatch (-want +got): \n%s", diff)
	}

	if _, err := json.Marshal(flatjson.DiffOperation{}); err == nil {
		t.Error("The diff operation type must be either add, remove, or replace.")
	}
}

func ExampleDiff() {
	diff, _ := flatjson.Diff(json1, json2)
	for _, op := range diff {
		switch op.Type {
		case flatjson.DiffOperationTypeAdd:
			fmt.Printf("%v %q: %v\n", op.Type, op.Path, op.Value)
		case flatjson.DiffOperationTypeRemove:
			fmt.Printf("%v %q\n", op.Type, op.Path)
		case flatjson.DiffOperationTypeReplace:
			fmt.Printf("%v %q: %v => %v\n", op.Type, op.Path, op.OldValue, op.Value)
		}
	}
	// Output:
	// replace []: map[address:map[city:Tokyo zip:100-0001] age:30 name:Alice] => map[address:map[city:Kyoto country:Japan] age:31 name:Alice]
	// replace ["address"]: map[city:Tokyo zip:100-0001] => map[city:Kyoto country:Japan]
	// replace ["address" "city"]: Tokyo => Kyoto
	// add ["address" "country"]: Japan
	// remove ["address" "zip"]
	// replace ["age"]: 30 => 31
}
