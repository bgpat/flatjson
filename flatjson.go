package flatjson

import (
	"cmp"
	"encoding/json"
	"slices"
)

// FlatJSON represents a flattened JSON structure as a slice of PathValue.
type FlatJSON []PathValue

// PathValue represents a single path-value pair in the flattened JSON structure.
type PathValue struct {
	Path  Path `json:"path"`
	Value any  `json:"value"`
}

// Flatten takes a nested JSON structure and returns flattened path-value pairs, sorted by the paths.
// The target can be any JSON-serializable object.
func Flatten(target any) (FlatJSON, error) {
	m, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}
	var u any
	if err := json.Unmarshal(m, &u); err != nil {
		// unreachable
		return nil, err
	}

	f := flatten(Path{}, u)
	slices.SortStableFunc(f, func(a, b PathValue) int {
		return cmp.Compare(a.Path.JSONPointer(), b.Path.JSONPointer())
	})
	return f, nil
}

func flatten(prefix Path, obj any) FlatJSON {
	f := FlatJSON{{
		Path:  prefix,
		Value: obj,
	}}

	if e, found := obj.(map[string]any); found {
		for p, v := range e {
			f = append(f, flatten(prefix.Join(p), v)...)
		}
	}

	if e, found := obj.([]any); found {
		for p, v := range e {
			f = append(f, flatten(prefix.Join(p), v)...)
		}
	}

	return f
}

// Get retrieves the value at the given path and returns it with whether it was found.
func (f FlatJSON) Get(p Path) (any, bool) {
	i, found := slices.BinarySearchFunc(f, p, func(a PathValue, b Path) int {
		return cmp.Compare(a.Path.JSONPointer(), b.JSONPointer())
	})
	if !found {
		return nil, false
	}
	return f[i].Value, true
}
