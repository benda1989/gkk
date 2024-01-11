package cache

import (
	"bytes"
	"net/http"
	"time"

	"github.com/benda1989/gkk/expect"
	"github.com/benda1989/gkk/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

// gin装饰器，参数：缓存时间
func Gin(expire time.Duration, keys ...func(*gin.Context) string) gin.HandlerFunc {
	return cacheGin(store, expire, keys...)
}
func cacheGin(store Store, expire time.Duration, keys ...func(*gin.Context) string) gin.HandlerFunc {
	sfGroup := singleflight.Group{}
	return func(c *gin.Context) {
		key := "cache:api:" + c.Request.RequestURI
		if keys != nil {
			key = "cache:api:" + keys[0](c)
		}
		resp := &Response{}
		if err := store.Get(key, &resp); err == nil {
			replyWithCache(c, resp)
			return
		} else {
			/////log
		}
		writer := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = writer
		raw, _, shared := sfGroup.Do(key, func() (any, error) {

			defer func(c *gin.Context) {
				if err := recover(); err != nil {
					if h, ok := err.(*expect.E); ok {
						h.Return(c)
					} else {
						c.Abort()
					}
				}
			}(c)
			c.Next()
			res := &Response{}
			res.set(writer)
			if !c.IsAborted() && c.Writer.Status() < 300 && c.Writer.Status() >= 200 {
				if err := store.Set(key, res, expire); err != nil {
					//////log
					logger.Log.Error(key, res)
				}
			}
			return res, nil
		})
		/////// 此处并发未暂测试,理论上同时访问的第二个起 shared即为true
		if shared {
			replyWithCache(c, raw.(*Response))
		}
	}
}

type Response struct {
	Status int
	Header http.Header
	Data   []byte
}

func (c *Response) set(g *responseWriter) {
	c.Status = g.Status()
	c.Data = g.body.Bytes()
	c.Header = g.Header().Clone()
}

type responseWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
func replyWithCache(c *gin.Context, respCache *Response) {
	c.Writer.WriteHeader(respCache.Status)
	for key, values := range respCache.Header {
		for _, val := range values {
			c.Writer.Header().Set(key, val)
		}
	}
	c.Writer.Write(respCache.Data)
	c.Abort()
}
