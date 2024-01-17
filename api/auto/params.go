package auto_gin

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"gkk/api"
	"gkk/code"
	"gkk/expect"
	"gkk/logger"
	"gkk/req"
	"gkk/str"

	"github.com/gin-gonic/gin"
)

var search = []string{">", ">=", "<", "<=", "like"}

func NewReqMap(obj any) *ReqMap {
	re := &ReqMap{
		tags:   map[string]string{},
		create: map[string]string{},
		update: map[string]string{},
		find:   map[string]findReq{},
		Order:  "id desc",
	}
	re.Encode(obj)
	re.fill(obj)
	return re
}

type ReqMap struct {
	tags     map[string]string
	create   map[string]string
	update   map[string]string
	find     map[string]findReq
	pk       string
	preload  []string
	Order    string
	OrderMap map[string]string
	source   reflect.Type
}

type findReq struct {
	Name      string
	Condition string
	IsMult    bool
}

func (r *ReqMap) Preload() []string {
	return r.preload
}
func (r *ReqMap) PK() string {
	return r.pk
}
func (r *ReqMap) Model() any {
	return reflect.New(r.source.Elem()).Interface()
}
func (r *ReqMap) Slice() any {
	sliceType := reflect.SliceOf(r.source)
	sliceVal := reflect.MakeSlice(sliceType, 0, 0)
	return sliceVal.Interface()
}
func (r *ReqMap) fill(obj any) {
	v := reflect.ValueOf(obj)
	r.source = reflect.TypeOf(obj)
	expect.PBM(v.Kind() != reflect.Ptr, "Use *Ptr")
	if n := v.MethodByName("OrderMap"); n.IsValid() {
		r.OrderMap = n.Interface().(func() map[string]string)()
	} else {
		r.OrderMap = map[string]string{
			"desc": r.pk + " desc",
			"asc":  r.pk,
		}
	}
}

func (r *ReqMap) splitTag(obj string) (c, u bool, s string) {
	c = true
	u = true
	for _, v := range strings.Split(obj, ";") {
		if v == "-" {
			c = false
			u = false
		} else if v == "unCreate" {
			c = false
		} else if v == "unUpdate" {
			u = false
		} else {
			s = v
		}
	}
	return
}

// cru标签规则：
// unCreate;unUpdate;like;=;>=;<=;>;<;in
func (r *ReqMap) Encode(obj any) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		cName := field.Tag.Get("cru")
		if cName == "-" {
			continue
		}
		isTime := field.Type.String() == "time.Time"
		if v.Field(i).Kind() == reflect.Struct && !isTime {
			r.Encode(v.Field(i).Interface())
			continue
		}
		name := str.ReplaceDownLine(field.Name)
		c, u, s := r.splitTag(cName)
		js := field.Tag.Get("json")
		if orm := field.Tag.Get("gorm"); orm != "-" {
			if js != "" && js != "-" {
				if orml := strings.ToLower(orm); !(strings.Contains(orml, "foreignkey") || strings.Contains(orml, "many2many")) {
					if strings.Contains(strings.ToLower(orm), "primarykey") {
						r.pk = js
					} else if c {
						r.create[js] = name
					}
					if u {
						r.update[js] = name
					}
				}
				if cName = field.Tag.Get("binding"); cName != "" {
					if isTime {
						// todo 前后端时间格式确定
						cName += ",datetime=2006-01-02 15:04"
					}
					r.tags[js] = cName
				} else if field.Type.String() != "string" {
					r.tags[js] = "numeric"
				} else if isTime {
					r.tags[js] = "datetime=" + time.RFC3339
				}
			}
			if orml := strings.ToLower(orm); strings.Contains(orml, "foreignkey") || strings.Contains(orml, "many2many") {
				r.preload = append(r.preload, field.Name)
			}
		}

		if len(s) > 0 {
			form := field.Tag.Get("form")
			if form == "" {
				if js == "" || js == "-" {
					logger.Log.Fatal(v.Type().String() + " " + field.Name + " 缺少form/json标签")
				}
				form = js
			}
			r.find[form] = findReq{name, s, s != "like" && len(s) > 3}
		}
	}
}
func (r *ReqMap) DecodeJson(c *gin.Context) map[string]any {
	reqs := map[string]any{}
	decoder := json.NewDecoder(c.Request.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&reqs); err != nil {
		expect.PMC("请检查json格式:"+err.Error(), code.PARAMETER_ERROR)
	}
	re := map[string]any{}
	data := r.create
	if c.Request.Method == "PUT" {
		data = r.update
	}
	for k, v := range data {
		if vv, ok := reqs[k]; ok {
			if va, ok := r.tags[k]; ok {
				api.DealValidError(req.Validate.Var(vv, va), "")
			}
			re[v] = vv
		} else {
			expect.PMC("参数遗漏: "+k, code.PARAMETER_ERROR)
		}
	}
	return re
}
func (r *ReqMap) DecodeFormWithPSO(c *gin.Context) (string, map[string]any, req.PageSizeOrder) {
	reqs := c.Request.URL.Query()
	re := map[string]any{}
	for k, v := range r.find {
		if val := reqs.Get(k); len(val) > 0 && val != "0" {
			if v.IsMult {
				re[v.Condition+" ?"] = val
			} else {
				if v.Condition == "in" {
					re[v.Name+" "+v.Condition+" ?"] = strings.Split(val, ",")
				} else {
					if v.Condition == "like" {
						if vals := strings.Split(val, ","); len(vals) > 1 {
							v.Condition += " '%" + strings.Join(vals[:len(vals)-1], "%' or "+v.Name+" like '%") + "%' or " + v.Name + " like "
							val = "%" + vals[len(vals)-1] + "%"
						} else {
							val = "%" + val + "%"
						}
					}
					re[v.Name+" "+v.Condition+" ?"] = val
				}
			}
		}
	}
	var pk string
	if val := reqs.Get(r.pk); len(val) > 0 && val != "0" {
		pk = val
	}
	pso := req.PageSizeOrder{Order: r.Order}
	pso.Bind(reqs)
	return pk, re, pso
}
func (r *ReqMap) DecodeForm(c *gin.Context) map[string]any {
	reqs := c.Request.URL.Query()
	re := map[string]any{}
	for k, v := range r.find {
		if val := reqs.Get(k); len(val) > 0 && val != "0" {
			if v.IsMult {
				re[v.Condition+" ?"] = val
			} else {
				if v.Condition == "in" {
					re[v.Name+" "+v.Condition+" ?"] = strings.Split(val, ",")
				} else {
					if v.Condition == "like" {
						if vals := strings.Split(val, ","); len(vals) > 1 {
							v.Condition += " '%" + strings.Join(vals[:len(vals)-1], "%' or "+v.Name+" like '%") + "%' or " + v.Name + " like "
							val = "%" + vals[len(vals)-1] + "%"
						} else {
							val = "%" + val + "%"
						}
					}
					re[v.Name+" "+v.Condition+" ?"] = val
				}
			}
		}
	}
	return re
}

func (r *ReqMap) CommentMap() map[string]string {
	re := map[string]string{}
	t := r.source.Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		cName := field.Tag.Get("gorm")
		if len(cName) > 0 && cName != "-" {
			for _, v := range strings.Split(cName, ";") {
				if strings.Contains(v, "comment") {
					if jName := field.Tag.Get("json"); len(cName) > 0 && cName != "-" {
						re[jName] = strings.Split(v, ":")[1]
					} else {
						re[str.ReplaceDownLine(field.Name)] = strings.Split(v, ":")[1]
					}
				}
			}
		}
	}
	return re
}
