package req

import (
	"net/url"
	"strconv"

	"gkk/code"
	"gkk/expect"
	"gkk/tool"

	"gorm.io/gorm"
)

type IdRequired struct {
	Id uint `form:"id" json:"id" binding:"required"`
}
type IdsReq struct {
	Id []uint `form:"ids" json:"ids" binding:"required"`
}
type IdReq struct {
	Id uint `gorm:"primaryKey" form:"id" json:"id"  cru:"unUpdate"`
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

func (p *PageSizeOrder) OffLimitOrder(db *gorm.DB, limit map[string]string, order string) *gorm.DB {
	db = p.PageSize.OffLimit(db)
	if p.Order != "" && tool.InArray(p.Order, limit) {
		db = db.Order(limit[p.Order])
	} else {
		db = db.Order(order)
	}
	return db
}
