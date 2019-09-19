package methods

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

var client fasthttp.Client
var subscribes = make(map[string]bool)

func getSubscribeID(db int, pattern string, callback string) string {
	return fmt.Sprintf("db:%d:pattern:%s:callback:%s", db, pattern, callback)
}