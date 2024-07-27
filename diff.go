package flatjson

import (
	"cmp"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// DiffOperation represents a difference operation between two JSON objects.
//
// Value is only set if the type is add or replace,
// and OldValue is only set if the type is replace.
type DiffOperation struct {
	Type     DiffOperationType `json:"type"`
	Path     Path              `json:"path"`
	Value    any               `json:"value"`
	OldValue any               `json:"old_value"`
}

// DiffOperationType defines the type of operations used in diffs.
type DiffOperationType string

// Operation types used in diffs.
// These are based on the `op` field in JSON Patch (RFC 6902).
const (
	DiffOperationTypeAdd     DiffOperationType = "add"
	DiffOperationTypeRemove  DiffOperationType = "remove"
	DiffOperationTypeReplace DiffOperationType = "replace"
)

// Diff computes the differences between this and another FlatJSON,
// returning operations that represent them.
func (f FlatJSON) Diff(x FlatJSON) []DiffOperation {
	var ops []DiffOperation

	for _, pv := range f {
		v, found := x.Get(pv.Path)
		if !found {
			ops = append(ops, DiffOperation{
				Type: DiffOperationTypeRemove,
				Path: pv.Path,
			})
			continue
		}

		if reflect.DeepEqual(pv.Value, v) {
			continue
		}
		ops = append(ops, DiffOperation{
			Type:     DiffOperationTypeReplace,
			Path:     pv.Path,
			Value:    v,
			OldValue: pv.Value,
		})
	}

	for _, pv := range x {
		if _, found := f.Get(pv.Path); found {
			continue
		}
		ops = append(ops, DiffOperation{
			Type:  DiffOperationTypeAdd,
			Path:  pv.Path,
			Value: pv.Value,
		})
	}

	slices.SortStableFunc(ops, func(a, b DiffOperation) int {
		return cmp.Compare(a.Path.JSONPointer(), b.Path.JSONPointer())
	})

	return ops
}

// Diff computes the differences between two JSON objects,
// returning operations that represent them.
//
// Errors may occur during marshaling and unmarshaling.
func Diff(x, y any) ([]DiffOperation, error) {
	xf, err := Flatten(x)
	if err != nil {
		return nil, err
	}

	yf, err := Flatten(y)
	if err != nil {
		return nil, err
	}

	return xf.Diff(yf), nil
}

// MarshalJSON implements encoding/json.Marshaler.
func (d DiffOperation) MarshalJSON() ([]byte, error) {
	switch d.Type {
	case DiffOperationTypeAdd:
		return json.Marshal(map[string]any{
			"type":  d.Type,
			"path":  d.Path,
			"value": d.Value,
		})
	case DiffOperationTypeRemove:
		return json.Marshal(map[string]any{
			"type": d.Type,
			"path": d.Path,
		})
	case DiffOperationTypeReplace:
		return json.Marshal(map[string]any{
			"type":      d.Type,
			"path":      d.Path,
			"value":     d.Value,
			"old_value": d.OldValue,
		})
	default:
		return nil, fmt.Errorf("unsupported diff operation type: %q", d.Type)
	}
}
