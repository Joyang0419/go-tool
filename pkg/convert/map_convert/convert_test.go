package map_convert

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
)

func TestKeys(t *testing.T) {
	type test struct {
		input map[string]int
		want  []string
	}
	tests := []test{
		{input: map[string]int{"one": 1, "two": 2}, want: []string{"one", "two"}},
		{input: map[string]int{}, want: []string{}},
		{input: nil, want: []string{}},
	}
	for _, tc := range tests {
		got := lo.Keys(tc.input)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Expected %v, but got %v", tc.want, got)
		}
	}
}

func TestValues(t *testing.T) {
	type test struct {
		input map[string]int
		want  []int
	}
	tests := []test{
		{input: map[string]int{"one": 1, "two": 2}, want: []int{1, 2}},
		{input: map[string]int{}, want: []int{}},
		{input: nil, want: []int{}},
	}
	for _, tc := range tests {
		got := lo.Values(tc.input)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Expected %v, but got %v", tc.want, got)
		}
	}
}
