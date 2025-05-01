package slice_convert

import (
	"fmt"
	"testing"
)

var testString = func(item string, index int) string {
	return fmt.Sprintf("%s %d", item, index)
}

var testInt = func(item int, index int) int {
	return item * index
}

func TestToExpectedTSlice(t *testing.T) {
	tt := []struct {
		name          string
		inputSlice    interface{}
		mapFn         interface{}
		expectedSlice interface{}
	}{
		{
			"String slice conversion",
			[]string{"one", "two", "three"},
			testString,
			[]string{"one 0", "two 1", "three 2"},
		},
		{
			"Int slice conversion",
			[]int{1, 2, 3},
			testInt,
			[]int{0, 2, 6},
		},
		{
			"Empty slice conversion",
			[]int{},
			testInt,
			[]int{},
		},
		{
			"Single item slice conversion",
			[]int{5},
			testInt,
			[]int{0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.inputSlice.(type) {
			case []string:
				actual := ToExpectedTSlice(tc.inputSlice.([]string), tc.mapFn.(func(item string, index int) string))
				expected := tc.expectedSlice.([]string)
				if !(len(actual) == len(expected)) || func() bool {
					for i, item := range actual {
						if item != expected[i] {
							return false
						}
					}
					return true
				}() {
					t.Fatalf("expected %v but got %v", expected, actual)
				}
			case []int:
				actual := ToExpectedTSlice(tc.inputSlice.([]int), tc.mapFn.(func(item int, index int) int))
				expected := tc.expectedSlice.([]int)
				if !(len(actual) == len(expected)) || func() bool {
					for i, item := range actual {
						if item != expected[i] {
							return false
						}
					}
					return true
				}() {
					t.Fatalf("expected %v but got %v", expected, actual)
				}
			}
		})
	}
}
