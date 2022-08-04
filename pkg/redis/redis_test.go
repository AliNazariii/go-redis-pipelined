package redis

import (
	"context"
	"testing"
	"time"
)

func TestSetAndPipelinedGet(t *testing.T) {
	redisModule := NewRedis()

	testData := map[string]string{
		"f123asda":      "2wdbnmqwd",
		"f12334sasda":   "2wdq1axwd",
		"f122sda3asda":  "2wd88qwd",
		"f12sd3asda":    "2wd11qwd",
		"f12dfsd63asda": "2wdd2qwd",
	}
	ctx := context.Background()

	for k, v := range testData {
		redisModule.Set(ctx, k, v, 2*time.Minute)
	}

	for k, v := range testData {
		result, err := redisModule.Get(ctx, k)
		if err != nil {
			t.Fail()
		}

		if result != v {
			t.Fail()
		}
	}

}
