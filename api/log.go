package api

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gkk/expect"
	"gkk/logger"
	"io"
	"net/url"
	"strings"
	"time"
)

func log(c *gin.Context) {
	// 开始时间
	startTime := time.Now()
	var body []byte
	if c.Request.Method != "GET" {
		if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			var e error
			if body, e = io.ReadAll(c.Request.Body); e == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}
	} // 处理请求
	c.Next()
	// 执行时间
	latencyTime := time.Now().Sub(startTime).String()
	// 日志格式
	path, _ := url.PathUnescape(c.Request.RequestURI)
	//
	field := map[string]any{
		"ip":     c.ClientIP(),
		"method": c.Request.Method,
		"url":    path,
		"code":   c.Writer.Status(),
		"period": latencyTime,
		"model":  "gin",
	}
	if len(logger.Log.Heads) > 0 {
		for _, v := range logger.Log.Heads {
			field[v] = c.Request.Header.Get(v)
		}
	}
	if err := c.Errors.Last(); err != nil {
		if h, ok := err.Meta.(*expect.E); ok {
			field["code"] = h.Code
			field["errMsg"] = h.Msg
			if h.Err != nil {
				field["err"] = h.Err.Error()
			}
		}
		logger.Log.WithFields(field).Error(string(body))
	} else {
		logger.Log.WithFields(field).Info(string(body))
	}
}
