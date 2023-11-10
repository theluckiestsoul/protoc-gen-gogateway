package generator

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func (gr *Generator) generateMiddlewares(svc *protogen.Service, g *protogen.GeneratedFile) {
	g.P("type unimplemented", svc.GoName, "Middleware struct {}")
	for _, method := range svc.Methods {
		g.P("func (u unimplemented", svc.GoName, "Middleware)", strings.ToLower(method.GoName), "Middlewares() []echo.MiddlewareFunc {")
		g.P("return []echo.MiddlewareFunc{}")
		g.P("}")
	}
	g.P()
	g.P("type ", strings.ToLower(svc.GoName), "Middleware struct {")
	g.P("unimplemented", svc.GoName, "Middleware")
	g.P("}")
	g.P()
}
