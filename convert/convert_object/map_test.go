package convert_object

import (
	"fmt"
	"testing"
)

func TestMapToStruct(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
		QQ   string
	}

	m := map[string]string{
		"Name": "John",
		"Age":  "30",
	}

	var s TestStruct
	s.QQ = "123456"
	MapToStruct(m, &s)

	fmt.Println(s.Age)
	fmt.Println(s.Name)
	fmt.Println(s.QQ)
}
