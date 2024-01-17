package cron

import (
	"reflect"
	"runtime"

	"gkk/middle"
	"gkk/tool"

	"github.com/robfig/cron"
)

var servers server

type server struct {
	Cron []string
	Func []func()
}

func (c *server) register(cron string, f func()) {
	c.Cron = append(c.Cron, cron)
	c.Func = append(c.Func, f)
}

func Run() {
	c := cron.New()
	for k, v := range servers.Cron {
		tool.Exit(c.AddFunc(v, wrap(servers.Func[k])), "添加定时器错误：")
	}
	c.Start()
}

func wrap(f func()) func() {
	return func() {
		defer middle.CronRecovery(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		f()
	}
}

// 参数：规则（s m h d m w），函数
func Register(rule string, f func()) {
	servers.register(rule, f)
}
