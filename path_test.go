package flatjson_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bgpat/flatjson"
)

func TestPathJoin(t *testing.T) {
	tests := []struct {
		prefix flatjson.Path
		name   any
		want   flatjson.Path
	}{
		{flatjson.Path{}, "key", flatjson.Path{"key"}},
		{flatjson.Path{"key1"}, "key2", flatjson.Path{"key1", "key2"}},
		{flatjson.Path{"key1"}, 0, flatjson.Path{"key1", 0}},
		{flatjson.Path{"key1", 0}, "key2", flatjson.Path{"key1", 0, "key2"}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v/%v", tt.prefix, tt.name), func(t *testing.T) {
			got := tt.prefix.Join(tt.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(%#v).Join(%q) = %q, want %q", tt.prefix, tt.name, got, tt.want)
			}
		})
	}
}

func TestPathJSONPointer(t *testing.T) {
	tests := []struct {
		path flatjson.Path
		want string
	}{
		{flatjson.Path{}, "/"},
		{flatjson.Path{"foo"}, "/foo"},
		{flatjson.Path{"foo", "bar"}, "/foo/bar"},
		{flatjson.Path{0}, "/0"},
		{flatjson.Path{"a/b"}, "/a~1b"},
		{flatjson.Path{"m~n"}, "/m~0n"},
		{flatjson.Path{"", "0"}, "//0"},
		{flatjson.Path{" "}, "/ "},
	}

	for _, tt := range tests {
		t.Run(tt.path.JSONPointer(), func(t *testing.T) {
			got := tt.path.JSONPointer()
			if got != tt.want {
				t.Errorf("(%#v).JSONPointer() = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func ExamplePath() {
	p := flatjson.Path{"store", "book", 0}
	title := p.Join("title")
	fmt.Printf("%#v\n", title)
	fmt.Printf("%q\n", title.JSONPointer())
	// Output:
	// flatjson.Path{"store", "book", 0, "title"}
	// "/store/book/0/title"
}
