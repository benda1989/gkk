package tool

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

const generateStr = "20060102150405"

func Struct2Map(req any) map[string]any {
	res := map[string]any{}
	if re, e := msgpack.Marshal(req); e != nil {
		fmt.Println("Struct2Map编码失败", e)
	} else {
		msgpack.Unmarshal(re, &res)
	}
	return res
}

func StructName(stt any) string {
	name := strings.Split(reflect.TypeOf(stt).String(), ".")
	return name[len(name)-1]
}

func IsEmpty(v any) bool {
	switch v.(type) {
	case time.Time:
		return v.(time.Time).IsZero()
	default:
		return IsBlank(reflect.ValueOf(v))
	}
}

func AppName() string {
	for i := 1; i < 10; i++ {
		pc, server, _, _ := runtime.Caller(i)
		if pc == 0 {
			return ""
		}
		if path.Base(server) == "main.go" {
			return path.Base(strings.Split(runtime.FuncForPC(pc).Name(), ".")[0])
		}
	}
	return ""
}

// 获取一个传入函数的名字
func FuncName(fun any) string {
	return runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
}

func FuncNameReal(fun any) string {
	res := strings.Split(FuncName(fun), ".")
	return res[len(res)-1]
}

// 获取层级调用的指定函数的上一层名字
func CallerFuncName(sts string) string {
	for i := 0; i < 10; i++ {
		pc, _, _, _ := runtime.Caller(i)
		if i > 0 && runtime.FuncForPC(pc).Name() == sts {
			pc, _, _, _ = runtime.Caller(i - 1)
			return runtime.FuncForPC(pc).Name()
		}
	}
	return ""
}

// 获取app的名字和main所在的路径
func CallerAppName() (string, string) {
	for i := 1; i < 10; i++ {
		pc, server, _, _ := runtime.Caller(i)
		if path.Base(server) == "main.go" {
			return path.Base(strings.Split(runtime.FuncForPC(pc).Name(), ".")[0]), server
		}
	}
	return "", ""
}

// 打印调用者的层级关系
func AllCaller() {
	for i := 1; i < 10; i++ {
		pc, server, _, _ := runtime.Caller(i)
		name := runtime.FuncForPC(pc).Name()
		fmt.Println(i, name, server)
	}
}

// reflect/value的数据填充
func FieldBlank(value, field reflect.Value) {
	switch value.Kind() {
	case reflect.String:
		if value.Len() == 0 && !IsBlank(field) {
			value.SetString(field.String())
		}
	case reflect.Bool:
		if !value.Bool() && !IsBlank(field) {
			value.SetBool(field.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == 0 && !IsBlank(field) {
			value.SetInt(field.Int())
		}
	default:

	}
}

// 判断reflect/value的是否为空
func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Chan:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func Key2Value[k comparable, v comparable](r map[k]v) map[v]k {
	re := make(map[v]k, len(r))
	for k, v := range r {
		re[v] = k
	}
	return re
}

func Reverse[T any](l []T) []T {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	return l
}

func InArray(find any, source any) bool {
	targetValue := reflect.ValueOf(source)
	switch reflect.TypeOf(source).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == find {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(find)).IsValid() {
			return true
		}
	}
	return false
}

func Duplicate[T comparable](req *[]T) (res []T) {
	tempMap := map[T]bool{}
	for _, e := range *req {
		if tempMap[e] == false {
			tempMap[e] = true
			res = append(res, e)
		}
	}
	return
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func UUID() string {
	fileNo, _ := uuid.NewRandom()
	return fileNo.String()
}

var num int64

func Generate(ns ...string) string {
	n := "D"
	if len(ns) > 0 {
		n = ns[0]
	}
	t := time.Now()
	s := t.Format(generateStr)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000000
	rs := sup(r, 3)
	return fmt.Sprintf(n+"%s%s%s%s", s[2:], ms, ps, rs)
}
func GenerateSimple(ns ...string) string {
	n := "D"
	if len(ns) > 0 {
		n = ns[0]
	}
	t := time.Now()
	s := t.Format(generateStr)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000000
	rs := sup(r, 3)
	return fmt.Sprintf(n+"%s%s%s", s[2:], ms, rs)
}

func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

func Exit(err any, str string) {
	switch arg := err.(type) {
	case error:
		if arg != nil {
			fmt.Println(str, arg.Error())
			os.Exit(1)
		}
	case bool:
		if arg {
			fmt.Println(str, err)
			os.Exit(1)
		}
	}
}
