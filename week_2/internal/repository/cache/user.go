package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kirby980/study/week_2/internal/domain"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// 传单机 Redis 可以
	// 传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

// A 用到了 B，B 一定是接口 => 这个是保证面向接口
// A 用到了 B，B 一定是 A 的字段 => 规避包变量、包方法，都非常缺乏扩展性
// A 用到了 B，A 绝对不初始化 B，而是外面注入 => 保持依赖注入(DI, Dependency Injection)和依赖反转(IOC)
// expiration 1s, 1m
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}
func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
func (cahce *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
