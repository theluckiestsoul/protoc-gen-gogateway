package generator

import (
	"fmt"
	"strings"

	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/models"
	"google.golang.org/protobuf/compiler/protogen"
)

func (gr *Generator) genUseCaseInterface(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("usecase/contract/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(fileName, "usecase")
	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package ", pkg)
	g.P()

	g.P("import (")
	g.P("\"context\"")
	g.P(")")

	for _, svc := range f.Services {
		g.P("type ", svc.GoName, "UseCase interface {")
		for _, method := range svc.Methods {
			g.P(method.GoName, "(ctx context.Context, req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ",error)")
		}
		g.P("}")
	}
}

func (gr *Generator) genUseCaseImplementation(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("usecase/implementation/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(fileName, "usecase")
	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package ", pkg)
	g.P()

	usecaseContractImport := fmt.Sprintf("%s/%s", usecaseContractPath, pkg)
	repositoryContractImport := fmt.Sprintf("%s/%s", repositoryContractPath, pkg)
	repoImplImport := fmt.Sprintf("%s/%s", repositoryImplPath, pkg)

	// Import repository
	g.P("import (")
	g.P("\"context\"")
	g.P("contract \"", usecaseContractImport, "\"")
	g.P("repoContract \"", repositoryContractImport, "\"")
	g.P("repoImpl \"", repoImplImport, "\"")
	g.P(")")

	for _, svc := range f.Services {
		//add unimplemented usecase
		svcName := strings.ToLower(svc.GoName)
		g.P("type ", svcName, "UnimplementedUseCase struct {")
		g.P("repo repoContract.", svc.GoName, "Repository")
		g.P("}")
		g.P()

		g.P("type ", svcName, "UseCase struct {")
		g.P("*", svcName, "UnimplementedUseCase")
		g.P("}")
		g.P()

		for _, method := range svc.Methods {
			g.P("func (s *", svcName, "UnimplementedUseCase) ", method.GoName, "(ctx context.Context, req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ",error) {")
			g.P("return s.repo.", method.GoName, "(ctx,req)")
			g.P("}")
		}
	}

	// generate new instance
	for _, svc := range f.Services {
		g.P("func New", svc.GoName, "UseCase() contract.", svc.GoName, "UseCase {")
		g.P("repo := repoImpl.New", svc.GoName, "Repository()")
		g.P("return &", strings.ToLower(svc.GoName), "UnimplementedUseCase{repo}")
		g.P("}")
	}
}
