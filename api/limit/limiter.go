package limit

import (
	"fmt"
	"net/http"
	"time"

	"gkk/captcha"
	"gkk/code"
	"gkk/expect"
	"gkk/logger"
	"gkk/req"

	"github.com/gin-gonic/gin"
)

// 参数：限制器（Newlimiter），可传参数：从gin的context提取key的方法
func Gin(limiter *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if limiter.KeyFun != nil {
			key = limiter.KeyFun(c)
		}
		if limiter.IPFun != nil && limiter.IPFun(key) {
			c.Next()
			return
		}
		l, err := limiter.Load(key)
		if err != nil {
			logger.Log.Error(err)
			c.Abort()
			return
		}
		//c.Header("X-Limit-Limit", strconv.FormatInt(l.Limit,10))
		//c.Header("X-Limit-Remaining", strconv.FormatInt(l.Remain,10))

		if l.Remain < 0 {
			if limiter.ValFun != nil {
				limiter.ValFun(key, c)
			} else {
				c.Header("X-Limit-Reset", l.Reset)
				c.JSON(http.StatusOK, expect.E{code.IP_LINMIT, "超过响应次数限制", nil})
			}
			c.Abort()
			return
		}
		c.Next()
	}
}

// 限制所有IP对当前路由的访问次数
func Path(expire time.Duration, limit int64, options ...LimitOption) gin.HandlerFunc {
	l := NewLimiter(expire, limit, options...)
	l.KeyFun = func(c *gin.Context) string {
		return c.Request.URL.Path
	}
	return Gin(l)
}

// 限制单个IP对当前路由的访问次数
func PathIP(expire time.Duration, limit int64, options ...LimitOption) gin.HandlerFunc {
	l := NewLimiter(expire, limit, options...)
	l.KeyFun = func(c *gin.Context) string {
		return c.ClientIP() + ":" + c.Request.URL.Path
	}
	return Gin(l)
}

// 限制单个IP对当前路由的访问次数并发送验证码
func PathIPCode(expire time.Duration, limit int64, options ...LimitOption) gin.HandlerFunc {
	l := NewLimiter(expire, limit, options...)
	if l.KeyFun == nil {
		l.KeyFun = func(c *gin.Context) string {
			return c.ClientIP() + ":" + c.Request.URL.Path
		}
	}
	l.ValFun = func(k string, c *gin.Context) {
		k = "code:" + k
		tmp := &captcha.CaptStore{Key: k}
		if err := l.Store.Get(k, &tmp); err != nil {
			value, item := captcha.DrawCaptcha()
			tmp.Value = value
			fmt.Println(tmp.Value)
			tmp.Base64 = item.EncodeB64()
			if err = l.Store.Set(k, tmp, time.Second*60); err != nil {
				c.JSON(http.StatusOK, expect.E{code.IP_LINMIT, "超过响应次数限制", nil})
				return
			}
		}
		c.JSON(200, req.Response{Code: code.IP_LINMIT_CODE, Data: tmp, Msg: "超过响应次数限制"})
	}
	return Gin(l)
}
