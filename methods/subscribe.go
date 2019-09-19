package methods

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"redis_key_observer/redis"
	"strconv"
	"time"
)

// http://localhost:8000/api/v1/subscribe?db=0&pattern=x*&interval=10&callback=http://localhost:8000/api/v1/test_callback
func Subscribe(ctx *fasthttp.RequestCtx, observer *redis.Observer) ([]byte, error) {
	db, err := strconv.Atoi(string(ctx.QueryArgs().Peek("db")))
	if err != nil || db < 0 || db > 9 {
		return nil, fmt.Errorf("param 'db': db < 0 || db > 9")
	}

	pattern := string(ctx.QueryArgs().Peek("pattern"))
	if pattern == "" {
		return nil, fmt.Errorf("param 'pattern': is empty")
	}

	interval, err := strconv.Atoi(string(ctx.QueryArgs().Peek("interval")))
	if err != nil || interval < 5 {
		return nil, fmt.Errorf("param 'interval': interval < 5")
	}

	callback := string(ctx.QueryArgs().Peek("callback"))
	if callback == "" {
		return nil, fmt.Errorf("param 'callback': is empty")
	}

	subscribeID := getSubscribeID(db, pattern, callback)
	if subscribes[subscribeID] {
		return nil, fmt.Errorf("already subscribed")
	}

	m := observer.Subscribe(db, pattern, time.Duration(interval)*time.Second, func(k, v string) {
		m := map[string]string{k: v}
		data, _ := json.Marshal(m)

		req := fasthttp.Request{}
		req.SetRequestURI(callback)
		req.SetBody(data)
		req.Header.SetMethod("POST")

		err := client.DoTimeout(&req, nil, 10*time.Second)
		if err != nil {
			log.Printf("callback: err=%v", err)
		}
	})

	subscribes[subscribeID] = true

	return json.Marshal(m)
}
