package server

import (
	"encoding/json"
	"github.com/osh88/redis_key_observer/redis"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type method struct {
	Name    string
	Handler Handler
}

type Handler func(*fasthttp.RequestCtx, *redis.Observer) ([]byte, error)

type server struct {
	httpServer *fasthttp.Server
	observer   *redis.Observer
	methods    map[string]method
}

func (s *server) Handler(ctx *fasthttp.RequestCtx) {
	method, ok := s.methods[string(ctx.Path())]
	if !ok {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	data, err := method.Handler(ctx, s.observer)
	if err != nil {
		data, _ = json.Marshal(struct{ Error string `json:"error"` }{Error: err.Error()})
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

func New(observer *redis.Observer) (*server, error) {
	s := &server{
		observer: observer,
		methods:  make(map[string]method),
	}

	return s, nil
}
