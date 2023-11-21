package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jellydator/ttlcache/v2"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"time"
)

type Store interface {
	Get(key string, value any) error
	Set(key string, value any, expire time.Duration) error
	Update(key string, value any) error
	Delete(key string) error
}

type Redis struct {
	RedisClient *redis.Client
}

func NewRedis(redisClient *redis.Client) *Redis {
	return &Redis{
		RedisClient: redisClient,
	}
}
func (store *Redis) Set(key string, value any, expire time.Duration) error {
	payload, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	return store.RedisClient.Set(ctx, key, payload, expire).Err()
}
func (store *Redis) Delete(key string) error {
	ctx := context.TODO()
	return store.RedisClient.Del(ctx, key).Err()
}
func (store *Redis) Get(key string, value any) error {
	ctx := context.TODO()
	if payload, err := store.RedisClient.Get(ctx, key).Bytes(); err == nil {
		return msgpack.Unmarshal(payload, value)
	} else {
		return err
	}
}
func (store *Redis) Update(key string, value any) error {
	return store.Set(key, value, -1)
}

type memory struct {
	Cache *ttlcache.Cache
}

func NewMemory(defaultExpiration time.Duration) *memory {
	cacheStore := ttlcache.NewCache()
	_ = cacheStore.SetTTL(defaultExpiration)
	// disable SkipTTLExtensionOnHit default
	cacheStore.SkipTTLExtensionOnHit(true)
	return &memory{
		Cache: cacheStore,
	}
}
func (c *memory) Set(key string, value any, expireDuration time.Duration) error {
	return c.Cache.SetWithTTL(key, value, expireDuration)
}
func (c *memory) Delete(key string) error {
	return c.Cache.Remove(key)
}
func (c *memory) Get(key string, value any) error {
	if val, err := c.Cache.Get(key); err != nil {
		return err
	} else {
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(val))
		return nil
	}
}
func (c *memory) GetWithTTL(key string, value any) (time.Duration, error) {
	if val, ttl, err := c.Cache.GetWithTTL(key); err != nil {
		return ttl, err
	} else {
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(val))
		return ttl, nil
	}
}
func (c *memory) Update(key string, value any) error {
	if _, ttl, err := c.Cache.GetWithTTL(key); err != nil {
		return c.Cache.SetWithTTL(key, value, ttl)
	} else {
		return err
	}
}
