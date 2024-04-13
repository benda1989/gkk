package gkk

import (
	"fmt"
	"github.com/benda1989/gkk/api"
	"github.com/benda1989/gkk/cache"
	"github.com/benda1989/gkk/config"
	"github.com/benda1989/gkk/cron"
	"github.com/benda1989/gkk/db"
	"github.com/benda1989/gkk/logger"
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
	// 定时服务
	cron.Run()
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
	if api.Router != nil {
		// gin服务
		api.Run()
	} else {
		select {}
	}
}
