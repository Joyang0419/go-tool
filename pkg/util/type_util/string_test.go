package type_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTString_Bytes(t *testing.T) {
	assert.Equal(t, []byte("test"), TString("test").Bytes())
}

func TestTString_String(t *testing.T) {
	assert.Equal(t, "test", TString("test").String())
}

func TestTString_Equals(t *testing.T) {
	s := TString("test")
	assert.True(t, s.Equals("test"))
}

func TestTString_IsEmpty(t *testing.T) {
	s := TString("")
	assert.True(t, s.IsEmpty())
}
