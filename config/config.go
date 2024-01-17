package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gkk/expect"
	"gkk/tool"

	"gopkg.in/yaml.v3"
)

var IsDebug bool

var c map[string]AppConf

func Init(f ...string) {
	var fp string
	if len(f) > 0 {
		fp = f[0]
	} else {
		projectPath, _ := os.Getwd()
		fp = filepath.Join(projectPath, "config.yaml")
	}
	_, err := os.Stat(fp)
	tool.Exit(err, "need config file default: "+fp)
	c = map[string]AppConf{}
	Read(fp, &c)
}

func getApp(app ...string) (re string) {
	if len(app) == 0 {
		re = tool.AppName()
	} else {
		re = app[0]
	}
	return
}

func setDefault(arg *AppConf) {
	d := c["default"]
	if arg.Qn == nil {
		arg.Qn = d.Qn
	}
	if arg.Queue == nil {
		arg.Queue = d.Queue
	}
	if arg.Cache == nil {
		arg.Cache = d.Cache
	}
	if arg.Db == nil {
		arg.Db = d.Db
	} else {
		obj := reflect.ValueOf(arg.Db).Elem()
		def := reflect.ValueOf(d.Db).Elem()
		for i := 0; i < obj.NumField(); i++ {
			tool.FieldBlank(obj.Field(i), def.Field(i)) //FieldByName(obj.Type().Field(i).Name))
		}
	}
}

func Read(path string, data any) {
	content, err := ioutil.ReadFile(path)
	expect.PEM(err, "读取配置失败")
	err = yaml.Unmarshal(content, data)
	expect.PEM(err, "映射配置失败")
}

func Write(path string, data any) {
	re, err := yaml.Marshal(data)
	expect.PEM(err, "映射配置失败")
	err = ioutil.WriteFile(path, re, 0777)
	expect.PEM(err, "保存配置失败")
}

func Get(app ...string) (re AppConf) {
	if value, ok := c[getApp(app...)]; ok {
		re = value
	} else {
		fmt.Println("缺少配置:", getApp(app...))
		os.Exit(1)
	}
	setDefault(&re)
	return
}
func GetDefault() AppConf {
	return c["default"]
}

func GetCustom(conf any, app ...string) {
	data, _ := yaml.Marshal(c[getApp(app...)].Custom)
	yaml.Unmarshal(data, conf)
}
