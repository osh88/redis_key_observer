package main

import (
	"flag"
	"fmt"
	"github.com/osh88/redis_key_observer/methods"
	"github.com/osh88/redis_key_observer/redis"
	"github.com/osh88/redis_key_observer/server"
	"github.com/valyala/fasthttp"
	"log"
	"runtime"
)

var (
	args struct {
		redisAddr string
		redisPass string
		httpAddr  string
	}
)

func init() {
	flag.StringVar(&args.redisAddr, "redisAddr", "", "Redis address")
	flag.StringVar(&args.redisPass, "redisPass", "", "Redis pass")
	flag.StringVar(&args.httpAddr, "httpAddr", ":8000", "HTTP address")
	flag.Parse()
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	o, err := redis.NewObserver(args.redisAddr, args.redisPass)
	check(err)

	s, err := server.New(o)
	check(err)

	s.SetHandler(`/api/v1/subscribe`, "Subscribe", methods.Subscribe)
	s.SetHandler(`/api/v1/unsubscribe`, "Unsubscribe", methods.Unsubscribe)
	s.SetHandler(`/api/v1/test_callback`, "testCallback", testCallback)

	log.Fatal(s.ListenAndServe(args.httpAddr))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func testCallback(ctx *fasthttp.RequestCtx, observer *redis.Observer) ([]byte, error) {
	fmt.Println("Method:", string(ctx.Method()))
	fmt.Println("Data:", string(ctx.Request.Body()))
	return nil, nil
}
