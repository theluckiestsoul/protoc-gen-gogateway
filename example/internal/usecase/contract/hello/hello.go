// Code generated by Go Gateway. DO NOT EDIT.

package hello

import (
	entity "github.com/theluckiestsoul/protoc-gen-gogateway/example/pkg/proto/entity"
)

import (
	"context"
)

type HelloServiceUseCase interface {
	SayHello(ctx context.Context, req *entity.HelloRequest) (*entity.HelloResponse, error)
	SayHelloAgain(ctx context.Context, req *entity.HelloRequest) (*entity.HelloResponse, error)
}
