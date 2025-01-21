package convert

import (
	"fmt"
	"testing"
)

func TestToJsonBytes(t *testing.T) {
	var TestStruct struct {
		Name string
		Age  int
		QQ   string
	}

	TestStruct.Name = "John"
	TestStruct.Age = 30
	TestStruct.QQ = "123456"

	jsonBytes, err := ToJsonBytes(TestStruct)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(jsonBytes))
}

func TestJsonBytesToObj(t *testing.T) {
	var TestStruct struct {
		Name string
		Age  int
		QQ   string
	}

	jsonBytes := []byte(`{"Name":"John","Age":30,"QQ":"123456"}`)

	err := JsonBytesToObj(jsonBytes, &TestStruct)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(TestStruct)
}
