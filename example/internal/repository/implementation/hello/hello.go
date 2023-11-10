// Code generated by Go Gateway. DO NOT EDIT. 

package hello

import (
	entity "github.com/theluckiestsoul/protoc-gen-gogateway/example/pkg/proto/entity"
)

import (
	contract "github.com/theluckiestsoul/protoc-gen-gogateway/example/internal/repository/contract/hello"
	"context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
)

type helloserviceRepository struct {
}

func (s *helloserviceRepository) SayHello(ctx context.Context, req *entity.HelloRequest) (*entity.HelloResponse, error) {
	url := "foo.googleapi.com:9090"
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	out := new(entity.HelloResponse)
	err = conn.Invoke(ctx, "/hello.HelloService/SayHello", req, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *helloserviceRepository) SayHelloAgain(ctx context.Context, req *entity.HelloRequest) (*entity.HelloResponse, error) {
	url := "foo.googleapi.com:9090"
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	out := new(entity.HelloResponse)
	err = conn.Invoke(ctx, "/hello.HelloService/SayHelloAgain", req, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewHelloServiceRepository() contract.HelloServiceRepository {
	return &helloserviceRepository{}
}
