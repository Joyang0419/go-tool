package convert_object

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
