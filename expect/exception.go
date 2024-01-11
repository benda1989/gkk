package expect

import (
	"fmt"
	"net/http"

	"github.com/benda1989/gkk/code"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type E struct {
	Code int   `json:"code"`
	Msg  any   `json:"msg"`
	Err  error `json:"-"`
}

func (r E) Error() string {
	if r.Err != nil {
		return fmt.Sprintf(`{\"code\":%d;\"msg\":\"%v\",\"err\":\"%s\"}`, r.Code, r.Msg, r.Err.Error())
	} else {
		return fmt.Sprintf(`{\"code\":%d;\"msg\":\"%v\"}`, r.Code, r.Msg)
	}
}
func (r E) Return(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func NE(msg any) error {
	return NEC(msg, code.ERROR)
}

func NEC(msg any, code int) error {
	return &E{
		Msg:  msg,
		Code: code,
	}
}

func NEE(err error, msg string) error {
	return NEEC(err, msg, code.ERROR)
}

func NEEC(err error, msg string, code int) error {
	if err != nil {
		return &E{
			Msg:  msg,
			Code: code,
			Err:  err,
		}
	}
	return nil
}

func PM(msg any) {
	panic(&E{code.ERROR, msg, nil})
}
func PMC(msg any, code int) {
	panic(&E{code, msg, nil})
}

func PBM(flag bool, msg any) {
	if flag {
		panic(&E{code.ERROR, msg, nil})
	}
}
func PBMC(flag bool, msg any, code int) {
	if flag {
		panic(&E{code, msg, nil})
	}
}

func PEM(err error, msg any) {
	if err != nil {
		panic(&E{code.ERROR, msg, err})
	}
}
func PEMC(err error, msg any, code int) {
	if err != nil {
		panic(&E{code, msg, err})
	}
}
func PDM(db *gorm.DB, msgs ...string) {
	if db.Error != nil {
		msg := "未查询到"
		if len(msgs) > 0 {
			msg = msgs[0]
		}
		panic(&E{code.ERROR, msg, db.Error})
	}
}
