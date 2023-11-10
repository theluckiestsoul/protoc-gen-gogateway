package generator

import (
	"fmt"
	"strings"

	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/models"
	"google.golang.org/protobuf/compiler/protogen"
)

var uniqueHandlers = make(map[string]handler)

type handler struct {
	services    []models.Service
	importAlias string
	importPath  string
}

// generate adapter
func (gr *Generator) generateAdapters(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("adapter/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(fileName, "adapter")
	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package ", pkg)
	g.P()

	handlerImportPath := fmt.Sprintf("%s/%s", handlerPath, pkg)
	usecaseImplImportPath := fmt.Sprintf("%s/%s", usecaseImplPath, pkg)

	// Import repository
	g.P("import (")
	g.P("handler \"", handlerImportPath, "\"")
	g.P("usecase \"", usecaseImplImportPath, "\"")
	g.P("echo\"", echov4Path, "\"")
	g.P(")")
	services := make([]models.Service, 0)
	for _, svc := range f.Services {
		if exists := gr.hasAnyHttpOption(svc); !exists {
			continue
		}
		services = append(services, models.Service{Name: "Register" + svc.GoName + "Handler"})
		g.P("func Register", svc.GoName, "Handler(g *echo.Group) {")
		g.P("h := handler.New" + svc.GoName + "Handler(usecase.New" + svc.GoName + "UseCase())")

		for _, method := range svc.Methods {
			option := gr.getHttpOption(method)
			if len(option.Path) == 0 || len(option.Method) == 0 {
				continue
			}

			//replace curly braces with : for path variable
			option.Path = strings.ReplaceAll(option.Path, "{", ":")
			option.Path = strings.ReplaceAll(option.Path, "}", "")

			if option.Method == "GET" {
				g.P("g.GET(\"", option.Path, "\", h.", method.GoName, ")")
			} else if option.Method == "POST" {
				g.P("g.POST(\"", option.Path, "\", h.", method.GoName, " )")
			} else if option.Method == "PUT" {
				g.P("g.PUT(\"", option.Path, "\", h.", method.GoName, ")")
			} else if option.Method == "DELETE" {
				g.P("g.DELETE(\"", option.Path, "\", h.", method.GoName, ")")
			} else if option.Method == "PATCH" {
				g.P("g.PATCH(\"", option.Path, "\", h.", method.GoName, ")")
			}
		}
		g.P("}")
	}
	if len(services) == 0 {
		g.Skip()
		return
	}
	if _, ok := uniqueHandlers[pkg]; !ok {
		uniqueHandlers[pkg] = handler{services: services, importAlias: pkg, importPath: fmt.Sprintf("%s/%s", adapterPath, pkg)}
	} else {
		handler := uniqueHandlers[pkg]
		handler.services = append(handler.services, services...)
		uniqueHandlers[pkg] = handler
	}
}
