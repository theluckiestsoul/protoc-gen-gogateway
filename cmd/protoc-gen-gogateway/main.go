package main

import (
	"flag"

	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/generator"
	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/models"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet
	mod := flags.String("mod", "", "module name")
	port := flags.String("port", "8080", "port to run server")

	options := protogen.Options{
		ParamFunc: flags.Set,
	}
	options.Run(func(gen *protogen.Plugin) error {
		gr := generator.NewGenerator(mod, port)
		if len(*mod) == 0 {
			panic("mod is required")
		}
		for i, f := range gen.Files {
			if !f.Generate {
				continue
			}
			file := models.NewFileInfo(f)
			//ignore if no service
			if len(file.Services) == 0 {
				continue
			}
			file.Validate()
			gr.GenerateFiles(gen, file)
			if len(gen.Files)-1 == i {
				gr.GenerateServer(gen, file)
			}
		}
		return nil
	})
}
