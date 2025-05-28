package server

import (
	"google.golang.org/grpc"
)

type IService interface {
	Register(*grpc.Server)
}
