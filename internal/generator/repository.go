package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/models"
	"google.golang.org/protobuf/compiler/protogen"
)

func (gr *Generator) genRepositoryInterface(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("repository/contract/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(fileName, "repository")
	g.P(GENERATOR_EDIT_WARNING)

	g.P("package ", pkg)

	repoMoq := make([]string, 0)
	repoMockGen := make([]string, 0)

	for _, svc := range f.Services {
		repoMoq = append(repoMoq, fmt.Sprintf("%sRepository:%sRepositoryMock", svc.GoName, svc.GoName))
		repoMockGen = append(repoMockGen, fmt.Sprintf("%sRepository=GoMock%sRepository", svc.GoName, svc.GoName))
	}

	g.P()

	g.P("import (")
	g.P("\"context\"")
	g.P(")")
	g.P()

	for _, svc := range f.Services {
		g.P("type ", svc.GoName, "Repository interface {")
		for _, method := range svc.Methods {
			g.P(method.GoName, "(ctx context.Context,req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ", error)")
		}
		g.P("}")
	}
}

func (gr *Generator) genRepositoryImplementation(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("repository/implementation/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(fileName, "repository")
	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package ", pkg)
	g.P()

	repoContractImport := fmt.Sprintf("%s/%s", repositoryContractPath, pkg)

	// generate import
	g.P("import (")
	g.P("contract\"", repoContractImport, "\"")
	g.P("\"context\"")
	g.P("grpc \"google.golang.org/grpc\"")
	g.P(strconv.Quote("google.golang.org/grpc/credentials"))
	g.P(strconv.Quote("crypto/tls"))
	g.P(")")
	g.P()
	for _, svc := range f.Services {
		svcName := strings.ToLower(svc.GoName)
		g.P("type ", svcName, "Repository struct {")
		g.P("}")
		g.P()
		for _, method := range svc.Methods {
			g.P("func (s *", svcName, "Repository) ", method.GoName, "(ctx context.Context, req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ",error) {")
			host := models.GetDefaultHost(svc)

			g.P("url:=", strconv.Quote(host+":"+*gr.port))
			g.P("conn,err:=grpc.Dial(url,grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))")
			g.P("if err!=nil{")
			g.P("return nil,err")
			g.P("}")
			g.P("defer conn.Close()")
			var endpoint string
			if f.File.Proto.Package != nil {
				endpoint = fmt.Sprintf("/%s.%s/%s", *f.File.Proto.Package, svc.Desc.Name(), method.Desc.Name())
			} else {
				endpoint = fmt.Sprintf("/%s/%s", svc.Desc.Name(), method.Desc.Name())
			}
			g.P("out := new(", method.Output.GoIdent, ")")
			g.P("err = conn.Invoke(ctx, ", strconv.Quote(endpoint), ", req, out)")
			g.P("if err != nil {")
			g.P("return nil, err")
			g.P("}")
			g.P("return out, nil")
			g.P("}")
			g.P()
		}
	}
	// generate new instance
	for _, svc := range f.Services {
		svcName := strings.ToLower(svc.GoName)
		g.P("func New", svc.GoName, "Repository() contract.", svc.GoName, "Repository {")
		g.P("return &", svcName, "Repository{}")
		g.P("}")
	}
}

func (gr *Generator) genMockRepository(gen *protogen.Plugin, f *models.FileInfo) {
	name, pkg := f.GetFileAndPkgName()
	fileName := fmt.Sprintf("repository/mock/%s/%s.go", pkg, name)
	g := gen.NewGeneratedFile(string(fileName), "repository")
	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package ", pkg)
	g.P()

	repoContractImport := fmt.Sprintf("%s/%s", repositoryContractPath, pkg)

	// generate import
	g.P("import (")
	g.P("contract\"", repoContractImport, "\"")
	g.P("faker\"", fakerPath, "\"")
	g.P("\"context\"")
	g.P(")")
	for _, svc := range f.Services {
		svcName := strings.ToLower(svc.GoName)
		g.P("type ", svcName, "Repository struct {")
		g.P("}")
		g.P()
		for _, method := range svc.Methods {
			g.P("func (s *", svcName, "Repository) ", method.GoName, "(ctx context.Context, req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ",error) {")
			g.P("res := &", method.Output.GoIdent, "{}")
			g.P("_ = faker.FakeData(res)")
			g.P("return res, nil")
			g.P("}")
		}
	}
	// generate new instance
	for _, svc := range f.Services {
		g.P("func New", svc.GoName, "Repository() contract.", svc.GoName, "Repository {")
		g.P("return &", strings.ToLower(svc.GoName), "Repository{}")
		g.P("}")
	}
}
