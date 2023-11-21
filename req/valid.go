package req

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// 定义一个全局翻译器T
var Trans ut.Translator
var Validate *validator.Validate

// InitTrans 初始化翻译器
func init() {
	var ok bool
	if Validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "_" {
				return ""
			}
			return name
		})
		uni := ut.New(en.New(), zh.New())
		Trans, _ = uni.GetTranslator("zh")
		err := zhTranslations.RegisterDefaultTranslations(Validate, Trans)

		Validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
			regRuler := "^1[3456789]{1}\\d{9}$"
			reg := regexp.MustCompile(regRuler)
			return reg.MatchString(fl.Field().String())
		})
		Validate.RegisterTranslation("mobile", Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "手机号错误!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
		if err != nil {
			fmt.Println("启动失败")
			os.Exit(1)
		}
	}
}
