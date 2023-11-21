package doc

import (
	"gkk/str"
	"reflect"
	"strings"
)

type structItem struct {
	Name     string
	Required bool
	Type     string
	Comment  string
}

func (i structItem) MD(need bool) string {
	if need {
		require := "否"
		if i.Required {
			require = "是"
		}
		return "|" + strings.Join([]string{i.Name, require, i.Type, i.Comment}, "|") + "|"
	}
	return "|" + strings.Join([]string{i.Name, i.Type, i.Comment}, "|") + "|"
}

func ParmaItem(obj any) []structItem {
	return paramItem(reflect.ValueOf(obj), "")
}

func paramItem(v reflect.Value, father string) []structItem {
	switch v.Kind() {
	case reflect.Map:
		return paramMap(v, father)
	case reflect.Interface, reflect.Ptr:
		v = v.Elem()
		return paramItem(v, father)
	case reflect.Slice, reflect.Array:
		if v.Len() > 1 {
			v = reflect.Indirect(v.Index(0))
			return paramItem(v, father)
		} else {

		}
	case reflect.Struct:
		if !strings.Contains(v.Type().String(), "time.Time") {
			return paramStruct(v, father)
		}
	}
	return nil
}

func paramMap(v reflect.Value, father string) (res []structItem) {
	if v.IsNil() {
		return paramStruct(v, v.Type().Key().String())
	}
	mi := v.MapRange()
	for i := 0; mi.Next(); i++ {
		value := mi.Value()
		if re := paramItem(value, str.ReplaceDownLine(str.String(mi.Key().Interface()))); re != nil {
			res = append(res, re...)
		} else {
			name := str.ReplaceDownLine(mi.Key().String())
			if father != "" {
				name = father + "." + name
			}
			if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			res = append(res, structItem{
				Name: name,
				Type: value.Kind().String(),
			})
		}
	}
	return
}

func paramStruct(v reflect.Value, father string) (res []structItem) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var js string
		tpName := field.Type.String()
		if field.Type.Kind() == reflect.Struct {
			if strings.Contains(strings.Split(tpName, ".")[1], field.Name) {
				js = father
				goto skip
			}
		}
		if js = field.Tag.Get("json"); js == "" {
			js = str.ReplaceDownLine(field.Name)
		}
		if father != "" {
			js = father + "." + js
		}
	skip:
		if re := paramItem(v.Field(i), js); re != nil {
			res = append(res, re...)
		} else {
			comment := field.Tag.Get("comment")
			if comment == "" {
				if cName := field.Tag.Get("gorm"); len(cName) > 0 {
					for _, v := range strings.Split(cName, ";") {
						if strings.Contains(v, "comment") {
							comment = strings.Split(v, ":")[1]
						}
					}
				}
			}
			res = append(res, structItem{
				Name:     js,
				Required: field.Tag.Get("binding") == "required",
				Type:     strings.Replace(tpName, "*", "", -1),
				Comment:  comment,
			})
		}
	}
	return
}
