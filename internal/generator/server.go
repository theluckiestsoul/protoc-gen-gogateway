package generator

import (
	"strconv"

	"github.com/theluckiestsoul/protoc-gen-gogateway/internal/models"
	"google.golang.org/protobuf/compiler/protogen"
)

func (gr *Generator) generateMain(gen *protogen.Plugin, f *models.FileInfo) {
	fileName := "adapter/main/main.go"
	g := gen.NewGeneratedFile(string(fileName), "adapter/server")

	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package main")
	g.P()

	g.P("import (")
	g.P("server \"", serverPath, "\"")
	g.P(")")
	g.P("func main() {")
	g.P("server.RunServer(\":" + *gr.port + "\")")

	g.P("}")
}

func (gr *Generator) generateServer(gen *protogen.Plugin, f *models.FileInfo) {

	fileName := "adapter/server/server.go"

	g := gen.NewGeneratedFile(string(fileName), "adapter/server")

	g.P(GENERATOR_EDIT_WARNING)
	g.P()
	g.P("package server")
	g.P()

	g.P("import (")
	g.P(strconv.Quote("context"))
	g.P(strconv.Quote("net/http"))
	g.P(strconv.Quote("os"))
	g.P(strconv.Quote("os/signal"))
	g.P(strconv.Quote("time"))
	for _, h := range uniqueHandlers {
		g.P(h.importAlias, " \"", h.importPath, "\"")
	}
	g.P("echo\"", echov4Path, "\"")
	g.P(")")
	g.P()

	g.P("type unimplementedMiddleware struct {}")
	g.P()
	g.P("func (u unimplementedMiddleware) getMiddlewares() []echo.MiddlewareFunc {")
	g.P("return []echo.MiddlewareFunc{}")
	g.P("}")
	g.P()
	g.P("type middleware struct {")
	g.P("unimplementedMiddleware")
	g.P("}")
	g.P()

	g.P("func RunServer(port string) {")
	g.P("e := echo.New()")
	g.P("g := e.Group(\"\")")
	g.P("m := middleware{}")
	g.P("g.Use(m.getMiddlewares()...)")
	g.P("e.GET(\"/\", welcome)")
	g.P("e.GET(\"/healthcheck\", healthCheck)")
	for _, h := range uniqueHandlers {
		for _, s := range h.services {
			g.P(h.importAlias, ".", s.Name, "(g)")
		}
	}
	g.P()
	g.P("e.Logger.Info(\"server started on port \", port)")
	g.P("// Start server")
	g.P("go func() {")
	g.P("if err := e.Start(port); err != nil && err != http.ErrServerClosed {")
	g.P("e.Logger.Fatal(\"shutting down the server\")")
	g.P("}")
	g.P("}()")
	g.P()
	g.P()

	g.P("// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.")
	g.P("// Use a buffered channel to avoid missing signals as recommended for signal.Notify")
	g.P("quit := make(chan os.Signal, 1)")
	g.P("signal.Notify(quit, os.Interrupt, os.Kill)")
	g.P("<-quit")
	g.P("ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)")
	g.P("defer cancel()")
	g.P("if err := e.Shutdown(ctx); err != nil {")
	g.P("e.Logger.Fatal(err)")
	g.P("}")
	g.P("}")
	g.P()
	g.P("func welcome(c echo.Context) error {")
	g.P("return c.String(http.StatusOK,", strconv.Quote("Go Gateway"), " )")

	g.P("}")

	g.P("func healthCheck(c echo.Context) error {")
	g.P("return c.JSON(http.StatusOK, \"OK\")")
	g.P("}")

}
