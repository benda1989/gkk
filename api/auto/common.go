package auto_gin

import (
	"gkk"
	"gkk/api"
	"gkk/compare"
	"gkk/expect"
	"gkk/req"
	"gkk/tool"
	"gorm.io/gorm"
)

type MG = map[string]GinHandler

func getDefaultMethod(key string) (string, string, func(*Context)) {
	if dC == nil {
		dC = new(defaultContext)
	}
	switch key {
	case "GET":
		return "GET", "", dC.Get
	case "POST":
		return "POST", "", dC.Post
	case "PUT":
		return "PUT", "", dC.Put
	case "DELETE":
		return "DELETE", "", dC.Delete
	case "ACTIVE":
		return "PUT", "/active", dC.Active
	case "SELECT":
		return "GET", "/select", dC.Select
	}
	gkk.Log.Fatal("Default 暂不支持该 " + key + " 方法")
	return "", "", nil
}

var dC *defaultContext

type defaultContext struct{}

func (c *defaultContext) Post(C *Context) {
	db, model := C.DBM()
	api.BindJson(C.Context, model)
	expect.PDM(db.Create(model), "创建失败")
	if C.auth != nil && !C.isAdmin {
		defaultDB.Model(model).Update(C.auth.Key(), C.auth.GetId())
	}
	if C.record != nil {
		if id, ok := tool.Struct2Map(model)[C.Req.PK()]; ok {
			C.record(C.auth, "POST", id, nil)
		} else {
			gkk.Log.Error("新建记录数据解析失败，请检查主键json配置是否一致")
		}
	}
}

func (c *defaultContext) Put(C *Context) {
	w := C.Json()
	id := w[C.Req.PK()]
	if C.auth != nil && !C.isAdmin {
		w[C.auth.Key()] = C.auth.GetId()
	}
	if C.record != nil {
		re := gkk.M{}
		defaultDB.Model(C.Req.Model()).First(re, id)
		C.record(C.auth, "PUT", id, compare.MapTrans(w, re, C.CommentMap()))
	}
	delete(w, C.Req.PK())
	expect.PDM(defaultDB.Model(C.Req.Model()).Where(C.Req.PK(), id).Updates(w), "更新失败")
}

func (c *defaultContext) Get(C *Context) {
	data, count := C.Get(nil, C.Req.Preload()...)
	C.RDS(data, count)
}

func (c *defaultContext) Delete(C *Context) {
	id := C.Id()
	if C.record != nil {
		re := gkk.M{}
		defaultDB.Model(C.Req.Model()).First(re, id)
		C.record(C.auth, "DELETE", id, nil)
	}
	if C.auth != nil && !C.isAdmin {
		defaultDB.Where(C.auth.Key(), C.auth.GetId()).Delete(C.Req.Model(), id)
	} else {
		defaultDB.Delete(C.Req.Model(), id)
	}
}

func (a *defaultContext) Active(c *Context) {
	id := c.Id()
	c.DB().Where(c.Req.PK(), id).Update("status", gorm.Expr("case when status = 1 then 2 else 1 end"))
	if c.record != nil {
		re := gkk.M{}
		defaultDB.Model(c.Req.Model()).First(re, id)
	}
}

func (a *defaultContext) Select(c *Context) {
	var res []req.IdName
	c.DB().Select("id,name").Scan(&res)
	c.RD(res)
}
