package memstore

import (
	"context"
	"time"

	"github.com/hugebear-io/gofiber/fabric"
	"github.com/redis/go-redis/v9"
)

type RedisMemstoreConfig struct {
	Address  string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type redisMemstore struct {
	rdb *redis.Client
}

func NewRedisMemstore(config *RedisMemstoreConfig) Memstore {
	opt := redis.Options{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	}

	rdb := redis.NewClient(&opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return redisMemstore{rdb: rdb}
}

func (m redisMemstore) Get(key string, val interface{}) error {
	raw, err := m.rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}

	if err := fabric.Recast(raw, val); err != nil {
		return err
	}
	return nil
}

func (m redisMemstore) Set(key string, val interface{}, duration *time.Duration) error {
	if err := m.rdb.Set(context.Background(), key, val, *duration).Err(); err != nil {
		return err
	}
	return nil
}

func (m redisMemstore) Delete(key string) error {
	if err := m.rdb.Del(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}
