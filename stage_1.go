package main

import "fmt"
import "github.com/go-redis/redis"
import "time"

type stage1 struct {
	Stage
}

func getStageOne() Stage {
	return Stage{
		name:    "Stage 1: PING <-> PONG",
		runFunc: runStage1,
	}
}

func runStage1(logger *customLogger) error {
	client := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		DialTimeout: 30 * time.Second,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		return err
	}

	if pong != "PONG" {
		return fmt.Errorf("Expected PONG, got %s", pong)
	}

	return nil
}
