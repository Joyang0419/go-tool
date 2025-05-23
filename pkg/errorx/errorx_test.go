package errorx

import (
	"context"
	"log/slog"
	"testing"

	"github.com/pkg/errors"

	"go-tool/pkg/slogx"
)

type Service struct{}

func (receiver *Service) Do() error {
	err := errors.New("error")
	// 我希望這邊印出的error是: errorx.Service.Do.error:

	return err
}

type Controller struct {
	Service *Service
}

func (receiver *Controller) Do() error {
	if err := receiver.Service.Do(); err != nil {
		return errors.Wrap(err, "Controller.Do")
	}
	return nil
}

func TestHasError(t *testing.T) {
	slogx.NewSlog(slogx.Config{
		CallerSkip: 7,
	})
	slog.SetDefault(
		slogx.NewSlog(slogx.Config{
			CallerSkip: 7,
		}),
	)
	s := new(Service)
	c := new(Controller)
	c.Service = s
	err := c.Do()
	ctx := context.Background()
	ctx = context.WithValue(ctx, slogx.TraceIDKey, "test-trace-id")
	LogError(ctx, err)
}
