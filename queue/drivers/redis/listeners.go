package drivers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func (rClient *RedisClient) StartConnection(wg *sync.WaitGroup) {
	go rClient.StartQueueListener(wg)
	go rClient.StartDelayedQueueListener(wg)
	go rClient.StartResultProcessor(wg)
	wg.Add(3)
}

func (rClient *RedisClient) StartQueueListener(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(rClient.Jobs)
	defer close(rClient.Results)
	
	log.Println("QUEUE LISTENER STARTED")
	
	for {
		msg, err := rClient.Client.RPop(ctx, rClient.Queue).Result()
		if err != nil && err != redis.Nil {
			log.Println("Error fetching immediate jobs:", err)
			continue
		}

		if msg != "" {
			rClient.Jobs <- msg
		}
	}
}

func (rClient *RedisClient) StartDelayedQueueListener(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("DELAYED JOBS QUEUE LISTENER STARTED")
	var delayedQueue = rClient.Queue + ":" + "delayed"
	
	for {
		now := time.Now().Unix()
		delayedJobs, err := rClient.Client.ZRangeByScore(ctx, delayedQueue, &redis.ZRangeBy{
			Min: "-inf",
			Max: fmt.Sprintf("%d", now),
		}).Result()

		if err != nil && err != redis.Nil {
			log.Fatal("Error fetching delayed jobs:", err)
			continue
		}

		for _, job := range delayedJobs {
			_, remErr := rClient.Client.ZRem(ctx, delayedQueue, job).Result()
			if remErr != nil {
				log.Fatal("Error removing job from delayed set:", remErr)
				continue
			}
			rClient.Jobs <- job
		}
	}
}

func (rClient *RedisClient) StartResultProcessor(wg *sync.WaitGroup) {
	defer wg.Done()
	for result := range rClient.Results {
		fmt.Println("Processed Result:", result)
	}
}