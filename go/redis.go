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

func (r RedisClient) SET(key string, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (r RedisClient) INCR(key string) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("INCR", key)
	return err
}

func (r RedisClient) GET(key string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("GET", key))
}

func (r RedisClient) MSET(m map[string][]byte) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("MSET", redis.Args{}.AddFlat(m)...)
	return err
}

func (r RedisClient) MGET(key []interface{}) ([][]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.ByteSlices(conn.Do("MGET", redis.Args{}.AddFlat(key)...))
}

func (r RedisClient) HSET(key, field string, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", key, field, value)
	return err
}

func (r RedisClient) HGET(key, field string, value interface{}) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("HGET", key, field))
}

func (r RedisClient) HMSET(key string, m map[string][]byte) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(m)...)
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
