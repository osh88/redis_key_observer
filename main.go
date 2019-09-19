package main

import (
	"flag"
	"fmt"
	"log"
	"redis_key_observer/methods"
	"redis_key_observer/redis"
	"redis_key_observer/server"
	"runtime"
	"time"
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

	o.Watch(0, "x*", 5*time.Second, func(k, v string) {
		fmt.Println(k, v)
	})

	time.Sleep(10 * time.Second)
	o.Unsubscribe(0, "x*")

	s, err := server.New()
	check(err)

	s.SetHandler(`/api/v1/watch`, "Watch", methods.Watch)

	log.Fatal(s.ListenAndServe(args.httpAddr))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
