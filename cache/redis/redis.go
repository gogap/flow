package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gogap/config"
	"github.com/gogap/flow/cache"
)

type RedisCache struct {
	redisClient *redis.Client
	keyExpir    time.Duration
	KeyPrefix   string
}

func init() {
	cache.RegisterCache("go-redis", NewRedisCache)
}

func NewRedisCache(opts ...cache.Option) (c cache.Cache, err error) {

	cacheOpts := cache.Options{}

	for _, o := range opts {
		o(&cacheOpts)
	}

	if cacheOpts.Config == nil {
		cacheOpts.Config = config.NewConfig()
	}

	redisCache := &RedisCache{
		redisClient: redis.NewClient(
			&redis.Options{
				Network:    cacheOpts.Config.GetString("network", "tcp"),
				Addr:       cacheOpts.Config.GetString("address", "127.0.0.1:6379"),
				Password:   cacheOpts.Config.GetString("password", ""),
				DB:         int(cacheOpts.Config.GetInt32("db", 0)),
				MaxRetries: int(cacheOpts.Config.GetInt32("max-retries", 0)),
				PoolSize:   int(cacheOpts.Config.GetInt32("pool-size", 10)),
			}),

		keyExpir:  cacheOpts.Config.GetTimeDuration("expiration", time.Minute*5),
		KeyPrefix: cacheOpts.Config.GetString("key-prefix"),
	}

	ret, err := redisCache.redisClient.Ping().Result()
	if err != nil {
		err = fmt.Errorf("connect redis failure: %s", err.Error())
		return
	}

	if ret != "PONG" {
		err = errors.New("ping redis failure")
		return
	}

	c = redisCache

	return
}

func (p *RedisCache) Set(k string, v interface{}) {
	p.redisClient.Set(p.KeyPrefix+k, v, p.keyExpir)
	return
}

func (p *RedisCache) Get(k string) (interface{}, bool) {
	ret, err := p.redisClient.Get(p.KeyPrefix + k).Result()

	if err != nil {
		return "", false
	}

	return ret, true
}

func (p *RedisCache) Delete(k string) {
	p.redisClient.Del(p.KeyPrefix + k)
}

func (p *RedisCache) Increment(k string, delta int64) (v int64, err error) {
	v, err = p.redisClient.IncrBy(p.KeyPrefix+k, delta).Result()
	return
}

func (p *RedisCache) Decrement(k string, delta int64) (v int64, err error) {
	v, err = p.redisClient.DecrBy(p.KeyPrefix+k, delta).Result()
	return
}

func (p *RedisCache) Flush() {
	p.redisClient.FlushDB()
}

func (p *RedisCache) IsLocal() bool {
	return false
}

func (p *RedisCache) CanStoreInterface() bool {
	return false
}
