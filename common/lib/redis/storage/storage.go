package storage

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/gomodule/redigo/redis"
)

var (
	RedisStorageClient *redis.Pool
)

func init() {
	redisHost, _ := config.String("storage_redis_host")
	redisPort, _ := config.Int("storage_redis_port")
	redisDb, _ := config.Int("storage_redis_db")
	address := fmt.Sprintf("%s:%d", redisHost, redisPort)

	RedisStorageClient = &redis.Pool{
		MaxIdle:     config.DefaultInt("storage_redis_maxidle", 4),
		MaxActive:   config.DefaultInt("storage_redis_maxactive", 512),
		IdleTimeout: 180 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", redisDb)
			return c, nil
		},
	}
}
