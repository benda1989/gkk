package gkk

import (
	"fmt"
	"reflect"

	"gkk/api"
	"gkk/cache"
	"gkk/config"
	"gkk/cron"
	"gkk/db"
	"gkk/logger"
)

var (
	Log *logger.Logger
)

type M = map[string]any
type MM = map[string]M
type MB = map[string]bool
type MS = map[string]string
type MSS = map[string][]string

func IsDebug() bool {
	return config.IsDebug
}

func init() {
	config.Init()
	Log = logger.Init()
	cache.Init()
}

func Run(apps ...any) {
	defer db.Close()
	// 注册用户app
	funCall(apps...)
	// 定时服务
	cron.Run()
	// gin服务
	api.Run()
	// 循环
	switch config.GetDefault().Db.User {

	default:
		fmt.Println(`
  ____    _    _   _    _
/  ___|  | | / /  | | / /
| |  _   | |/ /   | |/ /
| |_| |  | |\ \   | |\ \
 \____|  |_| \_\  |_| \_\

	`)
	}
	select {}
}

// funCall 函数在前参数 排后，reflect调用，传递多个函数。
func funCall(apps ...any) {
	all := make(map[reflect.Value][]reflect.Value)
	var flag reflect.Value
	for _, app := range apps {
		v := reflect.ValueOf(app)
		if v.Kind() == reflect.Func {
			flag = v
			all[v] = make([]reflect.Value, 0)
		} else {
			if len(all) == 0 {
				continue
			}
			all[flag] = append(all[flag], v)
		}
	}
	for app, arg := range all {
		app.Call(arg)
	}
}
