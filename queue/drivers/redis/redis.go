package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PratikKumar125/go-queue/queue/handler"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client  redis.Client
	Jobs    chan string
	Results chan string
	Queue string
	handler *handler.HandlerStruct
}

type JobDataStruct struct {
	Data map[string]any `json:"data"`
}

type RedisConnectionInputStruct struct {
	Addr string
	Password string
	DB int
	Queue string	
}

var ctx = context.Background()

func NewRedisQueue(inputs RedisConnectionInputStruct) *RedisClient {
	// creating the jobs store for the current queue
	handler := handler.NewHandler()
	handler.InitQueueJobs(inputs.Queue)
	fmt.Println("QUEUE JOBS STORE SETUP DONE")
	
	rdb := redis.NewClient(&redis.Options{
		Addr:     inputs.Addr,
		Password: inputs.Password,
		DB:       inputs.DB,
	})

	return &RedisClient{
		Client:  *rdb,
		Jobs:    make(chan string, 10),
		Results: make(chan string, 10),
		handler: handler,
		Queue: inputs.Queue,
	}
}

func (rClient *RedisClient) DispatchJob(inputs handler.DispatchJobStruct) (bool, error){
	if inputs.JobName == "" || inputs.Queue == "" {
		return false, fmt.Errorf("job name and queue are required")
	}
	if (inputs.Tries < 0) {
		return false, fmt.Errorf("tries cannot be negative")
	}

	inputs.Data["retryCount"] = 0
	done, err := rClient.AddJobToQueue(inputs.Queue, inputs)
	return done, err
}

func (rClient *RedisClient) AddJobToQueue(queue string, job handler.DispatchJobStruct) (bool, error) {
	var delayedQueue = queue + ":delayed"
	payload, err := json.Marshal(job)
	if err != nil {
		return false, fmt.Errorf("failed to marshal job data: %w", err)
	}

	// no delay provided
	if job.Delay <= -1 {
		_, err := rClient.Client.LPush(ctx, queue, payload).Result()
		if err != nil {
			return false, fmt.Errorf("failed to push job to queue: %w", err)
		}
		return true, nil
	}

	// delay provided
	log.Println("Adding in the delay queue", delayedQueue, job)
	delayedTimeInUnix := float64(time.Now().Add(time.Second * time.Duration(job.Delay)).Unix())
	_, err = rClient.Client.ZAdd(ctx, delayedQueue, redis.Z{
		Member: payload,
		Score:  delayedTimeInUnix,
	}).Result()
	if err != nil {
		return false, fmt.Errorf("failed to add job to delayed queue: %w", err)
	}

	return true, nil
}
