package inject

import (
	"testing"

	"go-tool/pkg/httpclient"
)

func TestFXModule(t *testing.T) {
	FXModule(httpclient.Config{})
}
