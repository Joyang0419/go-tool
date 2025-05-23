package errorx

import (
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func StackTrace(err error) errors.StackTrace {
	if st, ok := err.(stackTracer); ok {
		return st.StackTrace()
	}
	return nil
}
