package slice_convert

import (
	"github.com/samber/lo"
)

// ToExpectedTSlice converts inputSlice to expectedElement slice using mapFn.
func ToExpectedTSlice[inputElement, expectedElement any](
	inputSlice []inputElement,
	mapFn func(item inputElement, index int) expectedElement,
) []expectedElement {
	return lo.Map(inputSlice, mapFn)
}
