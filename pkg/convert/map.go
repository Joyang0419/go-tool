package convert

import (
	"github.com/mitchellh/mapstructure"
	pkgerrors "github.com/pkg/errors"
)

// MapToStruct 將 map 轉換為 struct; 仿照json.Unmarshal的概念
func MapToStruct[mapV any](m map[string]mapV, pointerVal any) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           pointerVal,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return pkgerrors.Cause(err)
	}

	return decoder.Decode(m)
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func MapKeysValues[K comparable, V any](m map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
