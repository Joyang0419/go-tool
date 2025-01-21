package convert

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

func TestMapKeys(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}

	keys := MapKeys(m)
	fmt.Println(keys)
}

func TestMapValues(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}

	values := MapValues(m)
	fmt.Println(values)
}
