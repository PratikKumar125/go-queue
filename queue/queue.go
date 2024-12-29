package queue

import (
	drivers "github.com/PratikKumar125/go-queue/queue/drivers/redis"
	"github.com/PratikKumar125/go-queue/queue/handler"
)

type (
	DispatchJobStruct = handler.DispatchJobStruct
	CreateJobStruct = handler.CreateJobStruct
	HandlerStruct     = handler.HandlerStruct
	RedisConnectionInputStruct = drivers.RedisConnectionInputStruct
)

var (
	NewHandler = handler.NewHandler
)

var (
	NewRedisQueue = drivers.NewRedisQueue
)

type (
	RedisClient = drivers.RedisClient
)
