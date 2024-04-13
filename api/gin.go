package api

import (
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/benda1989/gkk/api/doc"
	"github.com/benda1989/gkk/code"
	"github.com/benda1989/gkk/config"
	"github.com/benda1989/gkk/expect"
	"github.com/benda1989/gkk/middle"
	"github.com/benda1989/gkk/tool"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

const clear = `<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width, initial-scale=1"><title>微信缓存清理工具</title></head><body><button id='btn'>清理缓存</button><p id='p'></p><script>var len = localStorage.length;var arr = new Array();for(var i = 0; i < len; i++) {var getKey = localStorage.key(i);var getVal = localStorage.getItem(getKey);arr[i] = { 'key': getKey,'val': getVal};}const p = document.getElementById("p");if(arr.length>0){let aToStr=JSON.stringify(arr);p.innerHTML = aToStr;}else{p.innerHTML = "暂无缓存";}const btn = document.getElementById("btn");btn.onclick=function(){localStorage.clear();location.reload();}</script></body></html>`

type ginEngine struct {
	*gin.Engine
	port string
}

var Router *ginEngine

func New() *gin.Engine {
	c := config.Get().Gin
	if c.Mode == "debug" {
		config.IsDebug = true
	} else {
		gin.SetMode("release")
	}
	r := gin.New()
	if c.FrontDir != "" {
		r.Static("/index", c.FrontDir)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusForbidden, expect.E{Code: code.NOT_FOUND_ROUTE, Msg: "Not Found Route", Err: nil})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusForbidden, expect.E{code.NOT_FOUND_METH, "Not Found Method", nil})
	})
	r.Use(middle.Cors, log, middle.GinRecovery)

	if gin.IsDebugging() {
		doc.Init(r)
		r.GET("api/clear", func(c *gin.Context) {
			c.Status(200)
			writeContentType(c.Writer, htmlContentType)
			c.Writer.Write([]byte(clear))
		})
		r.POST("api/doc", handler)
		r.GET("api/doc", handler)
	}
	Router = &ginEngine{
		Engine: r,
		port:   c.Port,
	}
	return r
}

func Run() {
	s := endless.NewServer(":"+Router.port, Router.Engine)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	s.BeforeBegin = func(add string) {
		fmt.Printf("程序进程pid： %d\n", syscall.Getpid())
	}
	tool.Exit(s.ListenAndServe(), "启动失败:"+Router.port)
}

func handler(c *gin.Context) {
	if re := doc.Handler(c); re != nil {
		RD(c, re)
	}
}
