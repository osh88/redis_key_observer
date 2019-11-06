package methods

import (
	"fmt"
	"github.com/osh88/redis_key_observer/redis"
	"github.com/valyala/fasthttp"
	"strconv"
)

// http://localhost:8000/api/v1/unsubscribe?db=0&pattern=x*&callback=http://localhost:8000/api/v1/test_callback
func Unsubscribe(ctx *fasthttp.RequestCtx, observer *redis.Observer) ([]byte, error) {
	db, err := strconv.Atoi(string(ctx.QueryArgs().Peek("db")))
	if err != nil {
		return nil, fmt.Errorf("param 'db': %v", err)
	}
	if db < 0 || db > 9 {
		return nil, fmt.Errorf("param 'db': db < 0 || db > 9")
	}

	pattern := string(ctx.QueryArgs().Peek("pattern"))
	if pattern == "" {
		return nil, fmt.Errorf("param 'pattern': is empty")
	}

	callback := string(ctx.QueryArgs().Peek("callback"))
	if callback == "" {
		return nil, fmt.Errorf("param 'callback': is empty")
	}

	observer.Unsubscribe(db, pattern)

	delete(subscribes, getSubscribeID(db, pattern, callback))

	return nil, nil
}
