package captcha

import (
	"fmt"
	"time"

	"github.com/benda1989/gkk/cache"
	"github.com/benda1989/gkk/captcha"
	"github.com/benda1989/gkk/code"
	"github.com/benda1989/gkk/expect"
	"github.com/benda1989/gkk/tool"
	"github.com/gin-gonic/gin"
)

const key = "cache:captcha:"

func Create(app string) map[string]string {
	id := tool.Generate("CC")
	value, item := captcha.DrawCaptcha()
	if gin.IsDebugging() {
		fmt.Println(value)
	}
	expect.PBM(cache.Get().Set(key+":"+app+id, value, time.Second*60) != nil, "")
	return map[string]string{
		"key":    id,
		"base64": item.EncodeB64(),
	}
}

func Check(app, id, value string) {
	var c string
	if e := cache.Get().Get(key+":"+app+id, &c); e != nil {
		expect.PMC("验证码过期", code.CODE_EXPIRE)
	} else {
		if c != value {
			expect.PMC("验证码错误", code.CODE_WRONG)
		}
	}
}
