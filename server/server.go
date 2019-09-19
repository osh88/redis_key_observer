package server

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type method struct {
	Name    string
	Handler Handler
}

type Handler func(*fasthttp.RequestCtx) ([]byte, error)

type server struct {
	httpServer *fasthttp.Server
	methods    map[string]method
}

func (s *server) Handler(ctx *fasthttp.RequestCtx) {
	method, ok := s.methods[string(ctx.Path())]
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	data, err := method.Handler(ctx)
	if err != nil {
		data, _ = json.Marshal(struct {error string} {
			error: err.Error(),
		})
	}

	ctx.SetContentType(`application/json`)
	ctx.Response.Header.SetContentLength(len(data))
	ctx.Response.SetBody(data)
}

func (s *server) ListenAndServe(addr string) error {
	s.httpServer = &fasthttp.Server{
		Handler: s.Handler,
	}

	log.Printf("Запускаем сервис на %q\n", addr)
	return s.httpServer.ListenAndServe(addr)
}

func (s *server) SetHandler(path, name string, handler Handler) {
	s.methods[path] = method{Name: name, Handler: handler}
}

func New() (*server, error) {
	s := &server{
		methods: make(map[string]method),
	}

	return s, nil
}
