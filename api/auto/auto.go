package auto_gin

import (
	"github.com/gin-gonic/gin"
	"gkk"
	"gkk/api"
	"gkk/req"
	"gkk/tool"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var defaultDB *gorm.DB
var defaultAuth, defaultAdmin AuthHandler

type GinHandler func(*Context)
type AuthHandler func(*gin.Context) AuthInfoHandler

type AuthInfoHandler interface {
	Key() string
	GetId() any
	Info() (id any, nickname string, avatar string)
}

func Register(db *gorm.DB, auth, admin AuthHandler, g gin.IRoutes, models ...any) {
	BaseRegister(db, auth, admin)
	PathRegister(g, models...)
}

func BaseRegister(db *gorm.DB, args ...AuthHandler) {
	defaultDB = db
	if len(args) > 0 {
		defaultAuth = args[0]
		if len(args) > 1 {
			defaultAdmin = args[1]
		}
	}
}
func PathRegister(g gin.IRoutes, models ...any) {
	for _, model := range models {
		arg := NewReqMap(model)
		value := reflect.ValueOf(model)
		var recordFunc func(AuthInfoHandler, string, any, gkk.MSS)
		if r := value.MethodByName("Record"); r.IsValid() {
			record, ok := r.Interface().(func(AuthInfoHandler, string, any, gkk.MSS))
			if !ok {
				gkk.Log.Fatal(value.Type().String() + " Record 方法不符合func(auth AuthInfoHandler,method string, id any, args MSS) 请检查!")
			}
			recordFunc = record
		}
		def := value.MethodByName("Default")
		if def.IsValid() {
			if ginPaths, ok := def.Interface().(func() (string, string, string)); !ok {
				gkk.Log.Fatal(value.Type().String() + " Default 方法不符合func()(string, string, string) 请检查!")
			} else {
				orders, paths, methods := ginPaths()
				arg.Order = orders
				for _, v := range strings.Split(methods, ";") {
					if strings.Contains(paths, "admin") {
						m, p, f := getDefaultMethod(v)
						g.Handle(m, paths+p, newGinHandlerFunc(arg, f, defaultAdmin, true, recordFunc))
					} else {
						if strings.Contains(v, "Auth") {
							m, p, f := getDefaultMethod(strings.Replace(v, "Auth", "", -1))
							g.Handle(m, paths+p, newGinHandlerFunc(arg, f, defaultAuth, false, recordFunc))
						} else if strings.Contains(v, "Admin") {
							v = strings.Replace(v, "Admin", "", -1)
							m, p, f := getDefaultMethod(v)
							if strings.Count(methods, v) > 1 {
								g.Handle(m, paths+p+"/admin", newGinHandlerFunc(arg, f, defaultAdmin, true, recordFunc))
							} else {
								g.Handle(m, paths+p, newGinHandlerFunc(arg, f, defaultAdmin, true, recordFunc))
							}
						} else {
							m, p, f := getDefaultMethod(v)
							g.Handle(m, paths+p, newGinHandlerFunc(arg, f, nil, false, recordFunc))
						}
					}
				}
			}
		}
		paths := value.MethodByName("Paths")
		if !paths.IsValid() {
			continue
		}
		ginPaths, ok := paths.Interface().(func() MG)
		if !ok {
			gkk.Log.Fatal(value.Type().String() + " Paths 方法不符合func()gkk.MG 请检查!")
		}
		for p, f := range ginPaths() {
			vs := strings.Split(p, ":")
			if len(vs) <= 1 {
				gkk.Log.Fatal(value.Type().String() + " " + p + " 缺少参数：相对路径")
			}
			funName := tool.FuncNameReal(f)
			if strings.Contains(funName, "Auth") {
				g.Handle(vs[0], vs[1], newGinHandlerFunc(arg, f, defaultAuth, false, recordFunc))
			} else if strings.Contains(funName, "Admin") {
				g.Handle(vs[0], vs[1], newGinHandlerFunc(arg, f, defaultAdmin, true, recordFunc))
			} else {
				g.Handle(vs[0], vs[1], newGinHandlerFunc(arg, f, nil, false, recordFunc))
			}
		}
	}
}

func newGinHandlerFunc(r *ReqMap, f GinHandler, auth func(*gin.Context) AuthInfoHandler, isAdmin bool, record func(AuthInfoHandler, string, any, gkk.MSS)) gin.HandlerFunc {
	return func(c *gin.Context) {
		nc := &Context{
			Context: c,
			Req:     *r,
			isAdmin: isAdmin,
			record:  record,
		}
		if auth != nil {
			nc.auth = auth(c)
		}
		f(nc)
		if !nc.IfReturn {
			api.RM(c, "操作成功")
		}
	}
}

type Context struct {
	Context  *gin.Context
	Req      ReqMap
	auth     AuthInfoHandler
	record   func(AuthInfoHandler, string, any, gkk.MSS)
	isAdmin  bool
	IfReturn bool
}

func (C *Context) BindJson(data any) {
	api.BindJson(C.Context, data)
}
func (C *Context) BindParam(data any) {
	api.BindParam(C.Context, data)
}
func (C *Context) DBMS() (*gorm.DB, any) {
	return C.CheckAuth(defaultDB.Model(C.Req.Model())), C.Req.Slice()
}
func (C *Context) DBM() (*gorm.DB, any) {
	model := C.Req.Model()
	return C.CheckAuth(defaultDB.Model(model)), model
}
func (C *Context) DB() *gorm.DB {
	return C.CheckAuth(defaultDB.Model(C.Req.Model()))
}
func (C *Context) Id() uint {
	id := new(req.IdRequired)
	api.BindParam(C.Context, id)
	return id.Id
}
func (C *Context) Ids() []uint {
	id := new(req.IdsReq)
	api.BindParam(C.Context, id)
	return id.Id
}
func (C *Context) Form() gkk.M {
	return C.Req.DecodeForm(C.Context)
}
func (C *Context) FormPSO() (id string, where gkk.M, p req.PageSizeOrder) {
	return C.Req.DecodeFormWithPSO(C.Context)
}
func (C *Context) Json() gkk.M {
	return C.Req.DecodeJson(C.Context)
}
func (C *Context) JsonPK() (any, gkk.M) {
	re := C.Req.DecodeJson(C.Context)
	id := re[C.Req.PK()]
	delete(re, C.Req.PK())
	return id, re
}
func (C *Context) CommentMap() gkk.MS {
	return C.Req.CommentMap()
}
func (C *Context) ClearAuth() *Context {
	C.auth = nil
	return C
}
func (C *Context) CheckAuth(db *gorm.DB) *gorm.DB {
	if C.auth != nil && !C.isAdmin {
		db = db.Where(C.auth.Key(), C.auth.GetId())
	}
	return db
}
func (C *Context) Get(where gkk.M, preloads ...string) (any, int64) {
	id, w, p := C.FormPSO()
	if len(id) > 0 {
		w[C.Req.PK()] = id
		return C.GetOne(w, preloads...), -1
	} else {
		if where != nil {
			for k, v := range where {
				w[k] = v
			}
		}
		return C.List(w, p, preloads...)
	}
}
func (C *Context) GetOne(where gkk.M, preloads ...string) any {
	db, re := C.DBM()
	for k, v := range where {
		db = db.Where(k, v)
	}
	for _, v := range preloads {
		db = db.Preload(v)
	}
	db.Find(re)
	return re
}
func (C *Context) GetAll(where gkk.M, preloads ...string) any {
	db, res := C.DBMS()
	if where != nil {
		for k, v := range where {
			db = db.Where(k, v)
		}
	}
	for k, v := range C.Form() {
		db = db.Where(k, v)
	}
	for _, v := range preloads {
		db = db.Preload(v)
	}
	db.Order(C.Req.Order).Find(&res)
	return res
}

func (C *Context) List(w gkk.M, p req.PageSizeOrder, preloads ...string) (any, int64) {
	var count int64
	db, res := C.DBMS()
	for k, v := range w {
		db.Where(k, v)
	}
	db.Count(&count)
	for _, v := range preloads {
		db = db.Preload(v)
	}
	if C.Req.OrderMap != nil {
		p.OffLimitOrder(db, C.Req.OrderMap, C.Req.Order).Find(&res)
	} else {
		p.OffLimit(db).Order(C.Req.Order).Find(&res)
	}
	return res, count
}

func (C *Context) UserId() any {
	return C.auth.GetId()
}
func (C *Context) UserInfo() (id any, nickname string, avatar string) {
	return C.auth.Info()
}
func (C *Context) RD(data any) {
	C.IfReturn = true
	api.RD(C.Context, data)
}
func (C *Context) RDS(data any, count int64) {
	C.IfReturn = true
	if count < 0 {
		api.RD(C.Context, data)
	} else {
		api.RDS(C.Context, data, count)
	}
}
