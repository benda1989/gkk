package limit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gkk/api"
	"gkk/captcha"
	"gkk/code"
	"gkk/expect"
)

func Reset(c *gin.Context) {
	var s struct {
		Key string `json:"key" binding:"required"`
	}
	api.BindJson(c, &s)
	ss := captcha.CaptStore{Key: s.Key}
	if ss.Get() != nil {
		expect.PM("key不存在或已接口恢复")
	} else {
		value, item := captcha.DrawCaptcha()
		ss.Value = value
		fmt.Println(value)
		ss.Base64 = item.EncodeB64()
		if ss.Set() != nil {
			expect.PM("重置失败，稍后再试！")
		} else {
			api.RD(c, ss)
		}
	}
}
func Verify(c *gin.Context) {
	var s struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	api.BindJson(c, &s)
	ss := captcha.CaptStore{Key: s.Key}
	if ss.Get() != nil {
		ss.Delete()
		expect.PM("key不存在或已接口恢复")
	} else {
		if ss.Value == s.Value {
			ss.Delete()
		} else {
			expect.PMC("验证码错误", code.CODE_WRONG)
		}
	}
}
