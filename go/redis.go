package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisClient struct {
	pool *redis.Pool
}

var redisClient RedisClient

func NewRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	pool := &redis.Pool{
		MaxIdle:     180,
		MaxActive:   0,
		IdleTimeout: 30 * time.Second,
		Wait:        false,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port)) },
	}

	redisClient = RedisClient{
		pool: pool,
	}
}

func (r RedisClient) HMSET(key string, values ...interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", key, values)
	return err
}

func (r RedisClient) HMGET(key string, fields ...interface{}) ([][]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("HMGET", key, fields))
}

func (r RedisClient) FLUSH() {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Flush()
}
