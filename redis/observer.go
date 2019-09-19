package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
)

type Notify func(k, v string)

func NewObserver(addr, pass string) (*Observer, error) {
	options := &redis.Options{
		Addr:       addr,
		Password:   pass,
		MaxRetries: 5,
	}

	o := Observer{
		client:      redis.NewClient(options),
		unsubscribe: make(map[string]bool),
	}

	if err := o.client.ConfigSet("notify-keyspace-events", "AKE").Err(); err != nil {
		return nil, err
	}

	return &o, nil
}

type Observer struct {
	client      *redis.Client
	unsubscribe map[string]bool
	mu          sync.Mutex
}

// Подписывается на изменение строкового ключа в редисе
// db: номер базы редиса
// pattern: ключ в редисе
// interval: интервал проверок значения
// notify: функция обратного вызова при изменении ключа
func (o *Observer) Subscribe(db int, pattern string, interval time.Duration, notify Notify) map[string]string {
	values := make(map[string]string)
	for _, key := range o.client.Keys(pattern).Val() {
		values[key] = o.client.Get(key).Val()
	}

	prefix := fmt.Sprintf("__keyspace@%d__:", db)
	channel := prefix + pattern
	log.Println("Подписались на канал:", channel)

	go func() {
		sub := o.client.PSubscribe(channel)
		ch := sub.Channel()
		ticker := time.NewTicker(interval)
		check := func(key string) {
			v := o.client.Get(key).Val()
			if values[key] != v {
				values[key] = v
				notify(key, v)
			}
		}

		for !o.checkUnsubscribe(channel) {
			select {
			case msg := <-ch:
				key := msg.Channel[len(prefix):]
				check(key)
			case <-ticker.C:
				for _, key := range o.client.Keys(pattern).Val() {
					check(key)
				}
			}
		}

		if err := sub.PUnsubscribe(channel); err != nil {
			log.Println(err)
		}
		log.Println("Отписались от канала:", channel)
		o.unsubscribe[channel] = false
	}()

	return values
}

func (o *Observer) Unsubscribe(db int, pattern string) {
	channel := fmt.Sprintf("__keyspace@%d__:%s", db, pattern)
	o.mu.Lock()
	defer o.mu.Unlock()
	o.unsubscribe[channel] = true
}

func (o *Observer) checkUnsubscribe(channel string) bool {
	o.mu.Lock()
	o.mu.Unlock()
	return o.unsubscribe[channel]
}
