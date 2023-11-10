// Code generated by Go Gateway. DO NOT EDIT.

package hello

import (
	handler "github.com/theluckiestsoul/protoc-gen-gogateway/example/internal/handler/hello"
	usecase "github.com/theluckiestsoul/protoc-gen-gogateway/example/internal/usecase/implementation/hello"
	echo "github.com/labstack/echo/v4"
)

func RegisterHelloServiceHandler(g *echo.Group) {
	h := handler.NewHelloServiceHandler(usecase.NewHelloServiceUseCase())
	g.POST("/v1/messages:hello", h.SayHello)
	g.POST("/v1/messages:helloAgain", h.SayHelloAgain)
}
