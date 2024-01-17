package limit

import (
	"sync"
	"time"

	"gkk/cache"

	"github.com/gin-gonic/gin"
)

type LimitOption interface {
	apply(*Limiter)
}

type loption func(*Limiter)

func (o loption) apply(middleware *Limiter) {
	o(middleware)
}

// 限制器配置：达到限制处理
func Limited(f func(string, *gin.Context)) LimitOption {
	return loption(func(l *Limiter) {
		l.ValFun = f
	})
}

// 限制器配置：存储key的定制
func SetKey(f func(*gin.Context) string) LimitOption {
	return loption(func(l *Limiter) {
		l.KeyFun = f
	})
}

// 限制器配置：对key的放行/阻止
func SetIP(f func(string) bool) LimitOption {
	return loption(func(l *Limiter) {
		l.IPFun = f
	})
}

// 新建限制器，参数：限制周期，限制次数，自定义处理IP函数
func NewLimiter(expire time.Duration, limit int64, options ...LimitOption) *Limiter {
	l := &Limiter{
		Store:  cache.Get(),
		Period: expire,
		Limit:  limit,
	}
	for _, option := range options {
		option.apply(l)
	}
	return l
}

type limiterResult struct {
	Remain int64
	Reset  string
}

type Limiter struct {
	Store  cache.Store
	Period time.Duration
	Limit  int64
	IPFun  func(string) bool
	KeyFun func(*gin.Context) string
	ValFun func(string, *gin.Context)
}

func (l *Limiter) Load(key string) (*limiterResult, error) {
	key = "limit:" + key
	c := &counter{}
	if err := l.Store.Get(key, &c); err != nil {
		c.Limit = l.Limit
		c.Used = 1
		c.Expire = time.Now().Add(l.Period)
		if err = l.Store.Set(key, c, l.Period); err != nil {
			return &limiterResult{}, err
		}
	} else {
		if c.increase(1) {
			if err = l.Store.Update(key, c); err != nil {
				return &limiterResult{}, err
			}
		}
	}
	return c.result(), nil
}
func (l *Limiter) Reset(key string) {
	key = "limit:" + key
	l.Store.Delete(key)
}

type counter struct {
	Expire time.Time
	Limit  int64
	Used   int64
	mutex  sync.RWMutex
}

func (c *counter) increase(value int64) (re bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.Used <= c.Limit {
		c.Used += value
		re = true
	}
	return
}
func (c *counter) result() *limiterResult {
	return &limiterResult{
		c.Limit - c.Used,
		c.Expire.Format("2006-01-02 15:04:05"),
	}

}
func (c *counter) reset() {
	c.Used = 0
}
