package kv

import (
    "github.com/go-redis/redis/v7"
    "time"
)

type RedisDB struct {
    conn *redis.Client
}

func OpenRedis(address string, password string, DB int) (*RedisDB, error) {
    client := redis.NewClient(&redis.Options{
        Addr:address,
        Password:password,
        DB:DB,
    })
    _, err := client.Ping().Result()
    if err != nil {
        return nil, err
    }

    r := &RedisDB{
        conn:client,
    }
    return r, nil
}

func (r *RedisDB) Set(key string, value string, ttl time.Duration) error {
    return r.conn.Set(key, value, ttl).Err()
}

func (r *RedisDB) Get(key string) (string, error) {
    result, err := r.conn.Get(key).Result()
    return result, err
}
