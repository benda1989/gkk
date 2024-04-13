package api

import (
	"github.com/benda1989/gkk/db"
	"github.com/benda1989/gkk/expect"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
)

func NewModel[T any](order string, unUpdate ...string) *post[T] {
	re := post[T]{
		unPut: unUpdate,
	}
	re.get = NewFind[T](order)
	re.encodeReq(re.Model())
	return &re
}

type post[T any] struct {
	*get[T]
	unPut []string
	self  bool
}

func (r *post[T]) encodeReq(obj any) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		isTime := field.Type.String() == "time.Time"
		if v.Field(i).Kind() == reflect.Struct && !isTime {
			r.encode(v.Field(i).Interface())
			continue
		}
		cName := field.Tag.Get("json")
		if cName == "-" { //len(cName) <= 0 ||
			r.unPut = append(r.unPut, db.ColumnName(field.Name))
		}
	}
}

func (r *post[T]) copy(req any) (re *T) {
	r.DB().Where(req, r.pk).First(&re)
	return
}

func (r *post[T]) Req(c *gin.Context) *T {
	re := r.Model()
	BindJson(c, re)
	return re
}

func (r *post[T]) Create(data *T) *T {
	expect.PDM(r.DB().Omit(r.pk).Create(data), "创建失败")
	return data
}

func (r *post[T]) Update(new *T, w map[string]any) (*T, *T) {
	old := r.Model()
	tx := r.DB().Where(new, r.pk)
	for k, v := range w {
		tx = tx.Where(k, v)
	}
	expect.PDM(tx.First(old), "未查询到该数据")
	expect.PDM(r.DB().Where(old, r.pk).Omit(append(r.unPut, r.pk)...).Save(new), "更新失败")
	return new, old
}

func (r *post[T]) Active(id any, w map[string]any) {
	tx := r.DB().Where(r.pk, id)
	for k, v := range w {
		tx = tx.Where(k, v)
	}
	tx.Update("status", gorm.Expr("case when status = 1 then 2 else 1 end"))
	return
}

func (r *post[T]) Delete(ids []uint, where map[string]any) (res []*T) {
	tx := r.DB().Where(r.pk, ids)
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	tx.Find(&res)
	r.DB().Delete(&res)
	return
}
