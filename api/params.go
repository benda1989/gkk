package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gkk/api/doc"
	"gkk/code"
	"gkk/expect"
	"gkk/req"
	"strconv"
)

func Token(c *gin.Context) string {
	token := c.Request.Header.Get("Token")
	if token == "" {
		expect.PMC("此接口头部必须要传递token", code.NO_VALID_TOKEN)
	}
	return token
}
func AppCode(c *gin.Context) string {
	codes := c.Request.Header.Get("app")
	if codes == "" {
		expect.PMC("此接口头部必须要传递app", code.PARAMETER_ERROR)
	}
	return codes
}
func TokenBindJson(c *gin.Context, f any) string {
	BindJson(c, f)
	return Token(c)
}
func TokenBindParam(c *gin.Context, f any) string {
	BindParam(c, f)
	return Token(c)
}
func BindJson(c *gin.Context, f any) {
	err := c.ShouldBindJSON(f)
	if gin.IsDebugging() {
		c.Set("_req", f)
		doc.SetReq(c, f)
	}
	DealValidError(err, "JSON")
}
func BindJsonWith(c *gin.Context, f any) {
	DealValidError(c.ShouldBindBodyWith(f, binding.JSON), "JSON")
}
func BindParam(c *gin.Context, f any) {
	err := c.ShouldBindQuery(f)
	doc.SetReq(c, f)
	DealValidError(err, "FORM")
}

func DealValidError(err error, t string) {
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			expect.PM(RemoveTopStruct(err.(validator.ValidationErrors).Translate(req.Trans)))
		case *json.UnmarshalTypeError:
			expect.PMC(err.(*json.UnmarshalTypeError).Field+": 类型错误", code.PARAMETER_ERROR)
		case *strconv.NumError:
			expect.PMC(err.(*strconv.NumError).Num+": 类型错误", code.PARAMETER_ERROR)
		default:
			expect.PMC(fmt.Sprintf("请将参数放于%s中传递: ", t)+err.Error(), code.PARAMETER_ERROR)
		}
	}
}

func RemoveTopStruct(fields map[string]string) (re []string) {
	for _, err := range fields {
		re = append(re, err)
	}
	return
}
