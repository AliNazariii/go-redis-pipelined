package main

import (
	"context"
	"fmt"
	"redis-training/pkg/redis"
	"sync"
	"time"
)

func main() {
	redisModule := redis.NewRedis()

	ctx := context.Background()
	wg := &sync.WaitGroup{}
	testsSet := []string{"vaswd", "dvawwd", "dgdfg", "sd66", "asaz23"}
	for i, _ := range testsSet {
		wg.Add(1)
		go func(i int) {
			err := redisModule.Set(ctx, testsSet[i], "havij2", 2*time.Hour)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	testsGet := []string{"qwqw", "fgfvd", "cdaw43", "sfsegv5", "efwe21231"}
	for i, _ := range testsGet {
		wg.Add(1)
		go func(i int) {
			value, err := redisModule.Get(ctx, testsGet[i])
			if err != nil {
				fmt.Println("Error: ", err)
			}
			fmt.Println(value)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
