package api

import (
	"github.com/benda1989/gkk/db"
	"github.com/benda1989/gkk/logger"
	"github.com/benda1989/gkk/req"
	"github.com/benda1989/gkk/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var search = []string{">", ">=", "<", "<=", "=", "like", "in"}

func NewFind[T any](order string) *get[T] {
	r := &get[T]{
		order: order,
		find:  map[string]req.FindReq{},
	}
	r.encode(r.model)
	r.fill()
	return r
}

type get[T any] struct {
	model   T
	table   string
	order   string
	where   map[string]any
	pk      string
	preload []string
	find    map[string]req.FindReq
	comment map[string]string
}

func (r *get[T]) SetWhere(k string, v any) {
	r.where = map[string]any{k: v}
}

func (r *get[T]) Comment() map[string]string {
	if r.comment == nil {
		r.comment = map[string]string{}
		for _, v := range r.state().Schema.Fields {
			if len(v.Comment) > 0 {
				r.comment[v.Name] = v.Comment
			}
		}
	}
	return r.comment
}

func (r *get[T]) fill() {
	s := r.state()
	r.table = s.Table
	for _, v := range s.Schema.Fields {
		if v.PrimaryKey {
			r.pk = v.DBName
			break
		}
	}
	for k := range s.Schema.Relationships.Relations {
		r.preload = append(r.preload, k)
	}
}

func (r *get[T]) encode(obj any) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		cName := field.Tag.Get("cru")
		if len(cName) <= 0 || cName == "-" {
			continue
		}
		isTime := field.Type.String() == "time.Time"
		if v.Field(i).Kind() == reflect.Struct && !isTime {
			r.encode(v.Field(i).Interface())
			continue
		}
		if !tool.InArray(cName, search) {
			flage := true
			for _, v := range search {
				if strings.Contains(cName, v) {
					flage = false
					break
				}
			}
			if flage {
				logger.Log.Fatal(v.Type().String() + " " + field.Name + " 检查标签: " + cName)
			}
		}
		form := field.Tag.Get("form")
		if form == "" {
			form = field.Tag.Get("json")
			if form == "" || form == "-" {
				logger.Log.Fatal(v.Type().String() + " " + field.Name + " 缺少form/json标签")
			}
		}
		r.find[form] = req.FindReq{
			db.ColumnName(field.Name),
			cName,
			isTime,
		}
	}
}

func (r *get[T]) state() *gorm.Statement {
	db := db.DB.Model(r.model)
	db.Statement.Parse(r.model)
	return db.Statement
}

func (r *get[T]) DB() *gorm.DB {
	return db.DB.Model(r.model) //.Debug()
}

func (r *get[T]) Model() *T {
	re := r.model
	return &re
}

func (r *get[T]) Slice() []*T {
	return make([]*T, 0)
}

func (r *get[T]) First(w map[string]any, preloads ...string) *T {
	re := r.Model()
	db := db.W(w, r.DB())
	for _, v := range preloads {
		db = db.Preload(v)
	}
	db.First(&re)
	return re
}

func (r *get[T]) Find(w map[string]any, preloads ...string) (res []*T) {
	res = r.Slice()
	db := db.W(w, r.DB())
	for _, v := range preloads {
		db = db.Preload(v)
	}
	db.Order(r.order).Find(&res)
	return
}

func (r *get[T]) List(w map[string]any, p *req.PageSizeOrder, preloads ...string) (res []*T, count int64) {
	db := db.W(w, r.DB())
	db.Count(&count)
	if count > 0 {
		res = r.Slice()
		for _, v := range preloads {
			db = db.Preload(v)
		}
		p.OffLimitOrder(db, r.order).Find(&res)
	}
	return
}

func (r *get[T]) Select(w any, args ...any) (res []req.IdName) {
	r.DB().Where(w, args...).Select("id,name").Scan(&res)
	return
}

func (r *get[T]) Get(c *gin.Context) {
	id, w, p := r.FormPSO(c)
	for k, v := range r.where {
		w[k] = v
	}
	if len(id) > 0 {
		RD(c, r.First(w, r.preload...))
	} else {
		res, count := r.List(w, p, r.preload...)
		RDS(c, res, count)
	}
}

func (r *get[T]) FormPSO(c *gin.Context) (id string, re map[string]any, p *req.PageSizeOrder) {
	reqs := c.Request.URL.Query()
	re = map[string]any{}
	for k, v := range r.find {
		if val := reqs.Get(k); len(val) > 0 && val != "0" {
			fk, fv := v.Encode(val)
			re[fk] = fv
			if r.pk == v.Name {
				id = val
				return
			}
		}
	}
	p = &req.PageSizeOrder{}
	BindParam(c, p)
	return
}

func (r *get[T]) Form(c *gin.Context) map[string]any {
	reqs := c.Request.URL.Query()
	re := map[string]any{}
	for k, v := range r.find {
		if val := reqs.Get(k); len(val) > 0 && val != "0" {
			fk, fv := v.Encode(val)
			re[fk] = fv
		}
	}
	return re
}
