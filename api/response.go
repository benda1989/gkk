package api

import (
	"errors"
	"gkk/api/doc"
	"gkk/code"
	"gkk/expect"
	"gkk/js"
	"gkk/req"
	"net/http"

	"github.com/gin-gonic/gin"
)

var jsonContentType = []string{"application/json; charset=utf-8"}
var htmlContentType = []string{"text/html; charset=utf-8"}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	header["Content-Type"] = value
}

func response(data any, c *gin.Context) (re req.Response) {
	re.Msg = "操作成功"
	re.Data = data
	doc.SetRes(c, data)
	return
}

func RDE(c *gin.Context, data any, err error) {
	if err != nil {
		c.JSON(http.StatusOK, expect.E{Code: code.ERROR, Msg: "操作失败", Err: err})
	} else {
		RD(c, data)
	}
}

func RD(c *gin.Context, data any) {
	v, _ := js.Marshal(response(data, c))
	RDB(c, v)
}

func RDS(c *gin.Context, data any, count int64) {
	if data == nil {
		c.JSON(http.StatusOK, req.Response{Code: code.ERROR, Data: "数据为空", Msg: c.Request.Method + " " + c.Request.URL.String()})
	} else {
		v, _ := js.Marshal(response(req.List{
			List:  data,
			Total: count,
		}, c))
		RDB(c, v)
	}
}

func RDB(c *gin.Context, data []byte) {
	c.Status(200)
	writeContentType(c.Writer, jsonContentType)
	c.Writer.Write(data)
}

func RM(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, req.Response{Msg: msg})
}

func RE(c *gin.Context, err error) {
	RME(c, "操作成功", err)
}

func R(c *gin.Context) {
	RM(c, "操作成功")
}

func RME(c *gin.Context, msg string, err error) {
	if err != nil {
		var v *expect.E
		if errors.As(err, &v) {
			c.JSON(http.StatusOK, v)
		}
	} else {
		c.JSON(http.StatusOK, req.Response{Msg: msg})
	}
}
