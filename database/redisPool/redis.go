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
	if res != "OK" {
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


func (r *RedisPool) GetBin(key string) ([]byte, error) {
	conn := r.GetRedis()
	defer conn.Close()
	res, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}

func (r *RedisPool) HSet(key, filed, value string) error {
	conn := r.GetRedis()
	defer conn.Close()
	_, err := conn.Do("HSET", key, filed, value)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisPool) HGet(key, filed string) (string, error) {
	conn := r.GetRedis()
	defer conn.Close()
	return redis.String(conn.Do("HGET", key, filed))
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

func (r *RedisPool) ZADD(set string, key string, val interface{}) error {
	conn := r.GetRedis()
	defer conn.Close()
	_, err := conn.Do("ZADD", set, val, key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisPool) ZRANKE(set string, member string, rev bool) (int, error) {
	conn := r.GetRedis()
	defer conn.Close()
	op := ""
	if rev {
		op = "ZRANK"
	} else {
		op = "ZREVRANK"
	}
	res, err := redis.Int(conn.Do(op, set, member))
	if err != nil {
		return -1, err
	}
	return res, nil
}

type pair struct {
	member string
	score  string
}

func (r *RedisPool) ZRANGE(set string, start, end int, rev bool) (interface{}, error) {
	conn := r.GetRedis()
	defer conn.Close()

	op := ""
	if rev {
		op = "ZREVRANGE"
	} else {
		op = "ZRANGE"
	}
	var list []pair
	res, err := conn.Do(op, set, start, end, "withscores")
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("NOT OK")
	}
	islice := res.([]interface{})
	if len(islice)%2 != 0 {
		return nil, errors.New("error")
	}
	for i := 0; i < len(islice); i += 2 {
		member := string(islice[i].([]byte))
		score := string(islice[i+1].([]byte))
		list = append(list, pair{member: member, score: score})
	}

	return list, nil
}
