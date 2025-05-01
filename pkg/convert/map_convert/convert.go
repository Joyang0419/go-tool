package map_convert

import (
	"github.com/samber/lo"
)

func Keys[K comparable, V any](m map[K]V) []K {
	return lo.Keys(m)
}

func Values[K comparable, V any](m map[K]V) []V {
	return lo.Values(m)
}
