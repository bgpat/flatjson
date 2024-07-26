package flatjson

import (
	"fmt"
	"strconv"
	"strings"
)

// Path represents a sequence of elements that form a JSON path.
//
// The names in the Path are either string or int.
// A string represents the name of an object member, while an int represents an array index.
type Path []any

// Join adds a name and returns the new Path.
func (p Path) Join(name any) Path {
	r := make(Path, 0, len(p)+1)
	r = append(r, p...)
	r = append(r, name)
	return r
}

var jsonPointerReplacer = strings.NewReplacer(
	"~", "~0",
	"/", "~1",
)

// JSONPointer returns the JSON Pointer (RFC 6091) representation as a string.
func (p Path) JSONPointer() string {
	if len(p) == 0 {
		return "/"
	}

	var pointer string
	for _, i := range p {
		var s string
		switch v := i.(type) {
		case string:
			s = v
		case int:
			s = strconv.Itoa(v)
		default:
			// Fallback if the type is not string or int.
			s = fmt.Sprint(v)
		}
		pointer += "/" + jsonPointerReplacer.Replace(s)
	}
	return pointer
}
