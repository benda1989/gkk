package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gkk/code"
	"gkk/expect"
	"gkk/logger"
	"gkk/str"
	"net"
	"net/http"
	"os"
	"strings"
)

func CronRecovery(name string) {
	if err := recover(); err != nil {
		var codes int
		var msg any
		var er error
		if h, ok := err.(*expect.E); ok {
			codes = h.Code
			msg = h.Msg
			er = h.Err
		} else if e, ok := err.(error); ok {
			msg = fmt.Sprint("未知错误:", e.Error())
			logger.StackSend(5, e.Error())
			codes = code.UNKNOW_ERROR
		} else {
			msg = fmt.Sprint("服务器错误:", err)
			logger.StackSend(5, err.(string))
			codes = code.SERVER_ERROR
		}
		logger.Log.WithFields(map[string]any{
			"code":  codes,
			"model": "task",
			"func":  name,
			"err":   er,
		}).Error(msg)
	}
}

func GinRecovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			if brokenPipe {
				c.Error(err.(error))
			} else {
				if req, ok := c.Get("_req"); ok {
					logger.Log.Errorln("获取参数:", str.String(req))
				}
				if h, ok := err.(*expect.E); ok {
					h.Return(c)
					c.Errors = append(c.Errors, &gin.Error{Meta: h})
				} else if e, ok := err.(error); ok {
					logger.Log.Errorln("未知错误:", err)
					logger.StackSend(3, e.Error())
					c.JSON(http.StatusForbidden, expect.E{code.UNKNOW_ERROR, "未知错误", nil})
				} else {
					logger.Log.Errorln("服务器错误:", err)
					logger.StackSend(3, err.(string))
					c.JSON(http.StatusForbidden, expect.E{code.SERVER_ERROR, "服务器错误", nil})
				}
			}
			c.Abort()
		}
	}()
	c.Next()
}
