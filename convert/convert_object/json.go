package convert_object

import (
	"github.com/bytedance/sonic"
)

// ToJsonBytes 將物件轉換為 JSON 字串
func ToJsonBytes(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// JsonBytesToObj 將 JSON 字串轉換為指定物件
func JsonBytesToObj(jsonBytes []byte, val any) error {
	return sonic.Unmarshal(jsonBytes, val)
}
