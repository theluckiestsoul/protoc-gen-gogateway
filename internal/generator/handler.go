package generator

import (
	"fmt"

	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/models"
	"google.golang.org/protobuf/compiler/protogen"
)

func (gr *Generator) generateHandlers(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("handler/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(fileName, "handler")

	usecaseContractImport := fmt.Sprintf("%s/%s", usecaseContractPath, pkg)

	genImportStatement := false
	for _, svc := range f.Services {
		if exists := gr.hasAnyHttpOption(svc); !exists {
			continue
		}
		if !genImportStatement {
			g.P(GENERATOR_EDIT_WARNING)
			g.P()
			g.P("package ", pkg)
			g.P()
			g.P("import (")
			g.P("echo\"", echov4Path, "\"")
			g.P("contract \"", usecaseContractImport, "\"")
			if gr.isMapperRequired(svc) {
				g.P("mapper \"", mapperPath, "\"")
			}
			g.P("\"net/http\"")
			g.P(")")
			genImportStatement = true
		}

		g.P("type ", svc.GoName, "Handler struct {")
		g.P("usecase contract.", svc.GoName, "UseCase")
		g.P("}")
		g.P()
		g.P()
		g.P("func New", svc.GoName, "Handler(usecase contract.", svc.GoName, "UseCase) *", svc.GoName, "Handler {")
		g.P("return &", svc.GoName, "Handler{usecase}")
		g.P("}")
		g.P()
		g.P()
		for _, method := range svc.Methods {
			option := gr.getHttpOption(method)
			if len(option.Path) == 0 || len(option.Method) == 0 {
				continue
			}
			g.P("func (s *", svc.GoName, "Handler) ", method.GoName, "(c echo.Context) error {")
			g.P("req := &", method.Input.GoIdent, "{}")

			if option.Method == "GET" {
				//read from query
				g.P("if err := mapper.MapParams(c.QueryParams(), req, enum" + svc.Desc.FullName().Name() + "); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")

				//read from path
				g.P("if err := mapper.MapParams(mapper.AggregatePathParams(c.ParamNames(), c.ParamValues()), req, enum" + svc.Desc.FullName().Name() + "); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")
			} else if option.Method == "POST" {
				//read from body
				g.P("if err := c.Bind(req); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")
			} else if option.Method == "PUT" {
				//read from body
				g.P("if err := c.Bind(req); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")

				//read from path
				g.P("if err := mapper.MapParams(mapper.AggregatePathParams(c.ParamNames(), c.ParamValues()), req, enum" + svc.Desc.FullName().Name() + "); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")

			} else if option.Method == "DELETE" {
				//read from query
				g.P("if err := mapper.MapParams(c.QueryParams(), req, enum" + svc.Desc.FullName().Name() + "); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")

				//read from path
				g.P("if err := mapper.MapParams(mapper.AggregatePathParams(c.ParamNames(), c.ParamValues()), req, enum" + svc.Desc.FullName().Name() + "); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")
			} else if option.Method == "PATCH" {
				// read from body
				g.P("if err := c.Bind(req); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")

				//read from path
				g.P("if err := mapper.MapParams(mapper.AggregatePathParams(c.ParamNames(), c.ParamValues()), req, enum" + svc.Desc.FullName().Name() + "); err != nil {")
				g.P("return c.JSON(http.StatusBadRequest, err.Error())")
				g.P("}")
			}

			g.P("res,err := s.usecase.", method.GoName, "(c.Request().Context(), req)")
			g.P("if err != nil {")
			g.P("return err")
			g.P("}")
			g.P("return c.JSON(http.StatusOK, res)")
			g.P("}")
		}
	}
	if !genImportStatement {
		g.Skip()
	}
}
