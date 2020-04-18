package redisPool

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisPool struct {
	redisPool *redis.Pool
}

func NewDefaultRedisPool(pool *redis.Pool) *RedisPool {
	return &RedisPool{redisPool: pool}
}

// 获取redis连接
func (r *RedisPool) GetRedis() redis.Conn {
	return r.redisPool.Get()
}

func (r *RedisPool) CloseRedis() {
	if r.redisPool != nil {
		_ = r.redisPool.Close()
	}
}

func (r *RedisPool) Set(key string, val interface{}, ttl time.Duration) error {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := conn.Do("SET", key, val, "EX", ttl.Seconds())
	if err != nil {
		return err
	}
	if res != int64(1) {
		return errors.New("NOT OK")
	}
	return nil
}

func (r *RedisPool) Get(key string) (interface{}, error) {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := redis.Values(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *RedisPool) SetString(key, val string, ttl time.Duration) error {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := redis.String(conn.Do("SET", key, val, "EX", ttl.Seconds()))
	if err != nil {
		return err
	}
	if res != "OK" {
		return errors.New("NOT OK")
	}
	return nil
}

func (r *RedisPool) SetBin(key string, val []byte, ttl time.Duration) error {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := conn.Do("SET", key, val, "EX", ttl.Seconds())
	if err != nil {
		return err
	}
	if res != int64(1) {
		return errors.New("NOT OK")
	}
	return nil
}

// 获取缓存
func (r *RedisPool) GetString(key string) (string, error) {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RedisPool) GetUint64(key string) (uint64, error) {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := redis.Uint64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *RedisPool) SetUint64(key string, u uint64, ttl time.Duration) error {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := conn.Do("SET", key, u, "EX", ttl.Seconds())
	if err != nil {
		return err
	}
	if res != int64(1) {
		return errors.New("NOT OK")
	}
	return nil
}

func (r *RedisPool) GetBin(key string) ([]byte, error) {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}

func (r *RedisPool) Del(key string) error {
	conn := r.GetRedis()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	} else {
		return nil
	}
}
