package redis_wrapper

import "github.com/bilibili/kratos/pkg/cache/redis"

var (
    redisExecutor *redis.Redis
)

func Init(executor *redis.Redis)  {
    redisExecutor = executor
}

