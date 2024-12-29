package queue

import (
	drivers "github.com/PratikKumar125/go-queue/queue/drivers/redis"
	"github.com/PratikKumar125/go-queue/queue/handler"
)

type (
	DispatchJobStruct = handler.DispatchJobStruct
	CreateInputStruct = handler.CreateInputStruct
	HandlerStruct     = handler.HandlerStruct
)

var (
	NewHandler = handler.NewHandler
)

type (
	RedisClient = drivers.RedisClient
)

var (
	InitRedisQueue = drivers.InitRedisQueue
)
