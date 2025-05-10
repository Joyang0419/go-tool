package errorx

import (
	"github.com/pkg/errors"
)

func HasError(err error, allowedErrs ...error) bool {
	if err == nil {
		return false
	}

	for _, allowedErr := range allowedErrs {
		if errors.Is(err, allowedErr) {
			return false
		}
	}

	return true
}
