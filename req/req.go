package req

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/benda1989/gkk/code"
	"github.com/benda1989/gkk/expect"
	"gorm.io/gorm"
)

type IdRequired struct {
	Id uint `form:"id" json:"id" binding:"required"`
}
type IdsReq struct {
	Id []uint `form:"ids" json:"ids" binding:"required"`
}
type IdReq struct {
	Id uint `gorm:"primaryKey" form:"id" json:"id"  cru:"="`
}

func (a IdReq) Check() {
	expect.PBMC(a.Id <= 0, "参数必填:id", code.PARAMETER_ERROR)
}

type IdName struct {
	Id   uint   `json:"id" comment:"对应id"`
	Name string `json:"name" comment:"名称"`
}

type NameValue struct {
	Name  string `form:"name" json:"name"`
	Value any    `json:"value"`
}

type PageSize struct {
	Page int `json:"page" form:"page" cru:"-" comment:"第几页"`
	Size int `json:"size" form:"size" cru:"-" comment:"单叶几个"`
}

func (p *PageSize) load() (int, int) {
	if p.Size <= 0 {
		p.Size = 20
	}
	if p.Page > 1 {
		return p.Size, p.Size * (p.Page - 1)
	} else {
		return p.Size, 0
	}
}

func (p *PageSize) OffLimit(db *gorm.DB) *gorm.DB {
	l, o := p.load()
	return db.Offset(o).Limit(l)
}
func (p *PageSize) OffLimitDesc(db *gorm.DB, total int64) *gorm.DB {
	l, o := p.load()
	o = int(total) - o
	if o < l {
		l = o
		o = 0
	} else {
		o = o - l
	}
	return db.Offset(o).Limit(l)
}

func (p *PageSize) OffLimitS() (re string) {
	l, o := p.load()
	re = " limit " + strconv.Itoa(l)
	if l > 0 {
		re += " offset " + strconv.Itoa(o)
	}
	return
}

type PageSizeOrder struct {
	PageSize
	Order string `json:"order" form:"order" cru:"-"`
}

func (p *PageSizeOrder) Bind(form url.Values) {
	if v := form.Get("size"); len(v) > 0 {
		p.Size, _ = strconv.Atoi(v)
	}
	if v := form.Get("page"); len(v) > 0 {
		p.Page, _ = strconv.Atoi(v)
	}
	if order := form.Get("order"); len(order) > 0 {
		p.Order = order
	}
}

func (p *PageSizeOrder) OffLimitOrder(db *gorm.DB, order string) *gorm.DB {
	db = p.PageSize.OffLimit(db)
	if len(limitMap) > 0 {
		db = db.Order(limitMap[p.Order])
	} else {
		db = db.Order(order)
	}
	return db
}

var limitMap = map[string]string{}

func RegisterLimitMap(m map[string]string) {
	for k, v := range m {
		limitMap[k] = v
	}
}

type FindReq struct {
	Name      string
	Condition string
	Iftime    bool
}

func (v *FindReq) time(val string) any {
	if v.Iftime {
		re, err := time.Parse("2006-01-02T15:04:05", val)
		expect.PEM(err, v.Name+" 请使用格式：2006-01-02T15:04:05")
		return re
	} else {
		return val
	}
}

func (v *FindReq) Encode(val string) (string, any) {
	switch v.Condition {
	case "like":
		if vals := strings.Split(val, ","); len(vals) > 1 {
			v.Condition += " '%" + strings.Join(vals[:len(vals)-1], "%' or "+v.Name+" like '%") + "%' or " + v.Name + " like "
			val = "%" + vals[len(vals)-1] + "%"
		} else {
			val = "%" + val + "%"
		}
		fallthrough
	case ">", ">=", "<", "<=", "=":
		return v.Name + " " + v.Condition + " ?", v.time(val)
	case "in":
		return v.Name + " " + v.Condition + " ?", strings.Split(val, ",")
	default:
		return v.Condition + " ?", v.time(val)
	}
}
