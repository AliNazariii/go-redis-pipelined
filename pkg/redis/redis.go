package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
}

type GetResponse struct {
	Value interface{}
	Error error
}

type Impl struct {
	Client   *redis.Client
	getQueue map[string][]chan GetResponse
	getLock  sync.Mutex
}

func NewRedis() Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:8181",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	getQueue := make(map[string][]chan GetResponse)
	redisModule := &Impl{
		Client:   rdb,
		getQueue: getQueue,
	}
	redisModule.BatchGetter()
	return redisModule
}

func (r *Impl) BatchGetter() {
	ticker := time.NewTicker(3 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				r.getLock.Lock()
				oldQueue := r.getQueue
				r.getQueue = make(map[string][]chan GetResponse)
				r.getLock.Unlock()

				pipe := r.Client.Pipeline()
				ctx := context.Background()

				temp := make(map[string]*redis.StringCmd, 0)
				for k, _ := range oldQueue {
					t := pipe.Get(ctx, k)
					temp[k] = t
				}

				_, err := pipe.Exec(ctx)
				if err != nil {
					fmt.Println("error in batch getter", err)
				}

				for key, _ := range oldQueue {
					get, err := temp[key].Result()
					for c, _ := range oldQueue[key] {
						oldQueue[key][c] <- GetResponse{
							Value: get,
							Error: err,
						}
					}
				}
			}
		}
	}()
}

func (r *Impl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

//func (r *Impl) Get(ctx context.Context, key string) (interface{}, error) {
//	value, err := r.Client.Get(ctx, key).Result()
//	return value, err
//}

func (r *Impl) Get(ctx context.Context, key string) (interface{}, error) {
	keyChannel := make(chan GetResponse)
	r.getLock.Lock()
	queue, ok := r.getQueue[key]
	if ok && len(queue) == 0 {
		r.getQueue[key] = make([]chan GetResponse, 0)
	}
	r.getQueue[key] = append(r.getQueue[key], keyChannel)
	r.getLock.Unlock()
	result := <-keyChannel
	return result.Value, result.Error
}
