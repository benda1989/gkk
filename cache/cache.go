package cache

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/benda1989/gkk/config"
	"github.com/benda1989/gkk/logger"
	"github.com/benda1989/gkk/tool"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
)

var store Store

func Get() Store {
	if store == nil {
		fmt.Println(" please add cache.Init()")
		os.Exit(1)
	}
	return store
}

func Init() {
	c := config.GetDefault().Cache
	if c != nil && c.Host != "" {
		fmt.Println("缓存配置： using Redis: " + c.Host)
		store = NewRedis(redis.NewClient(&redis.Options{
			Addr:     c.Host,
			DB:       c.Db,
			Password: c.Password,
		}))
		return
	}
	fmt.Println("缓存配置： using MEM")
	store = NewMemory(1 * time.Minute)
}

// 参数：缓存时间，函数输入，结果接受(指针)，失效调用函数
func Cache(expire time.Duration, input, output, fun any) {
	cache(store, expire, input, output, fun)
}

func cache(store Store, expire time.Duration, input, output, fun any) {
	name := tool.FuncName(fun)
	key := "cache:func:" + name + ":" + tool.MD5Struct(input)
	if err := store.Get(key, output); err != nil {
		f := reflect.ValueOf(fun)
		if f.Kind() != reflect.Func {
			logger.Log.Fatal(name + ":非可执行函数")
		}
		o := reflect.ValueOf(output)
		if o.Kind() != reflect.Ptr {
			logger.Log.Fatal(name + ":output必须是指针类型")
		}
		sfGroup := singleflight.Group{}
		raw, _, shared := sfGroup.Do(key, func() (any, error) {
			defer func() {
				if e := recover(); e != nil {
					logger.Log.Error("Cache", e)
				}
			}()
			f.Call([]reflect.Value{reflect.ValueOf(input), o})
			if err = store.Set(key, o.Elem().Interface(), expire); err != nil {
				//////log
				logger.Log.Error("Cache", input, output)
			}
			return output, nil
		})
		if shared {
			output = raw
		}
	}
}
