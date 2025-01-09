package grpcx_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/processout/grpc-go-pool"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type PoolParams struct {
	ClientTimeout   time.Duration
	URL             string
	InitSize        int
	PoolSize        int
	IdleDuration    time.Duration
	MaxLifeDuration time.Duration
	UseTLS          bool
}

func NewPool(params PoolParams, opt ...grpc.DialOption) *grpcpool.Pool {
	var options []grpc.DialOption
	if params.UseTLS {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options, opt...)
	factory := func() (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), params.ClientTimeout)
		defer cancel()

		return grpc.DialContext(ctx, params.URL, options...)
	}

	var (
		pool *grpcpool.Pool
		err  error
	)
	if params.MaxLifeDuration == 0 {
		pool, err = grpcpool.New(factory, params.InitSize, params.PoolSize, params.IdleDuration)
	} else {
		pool, err = grpcpool.New(factory, params.InitSize, params.PoolSize, params.IdleDuration, params.MaxLifeDuration)
	}

	if err != nil {
		panic(fmt.Sprintf("[NewPool]grpcpool.New error: %v", err))
	}

	return pool
}

func InjectPool(params PoolParams) fx.Option {
	return fx.Options(
		fx.Supply(params),
		fx.Provide(NewPool),
	)
}
