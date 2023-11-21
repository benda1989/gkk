package doc

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type apiItems []*apiItem
type cateApiItems []*cateApiItem

func (a *apiItems) setReq(c *gin.Context, arg any) {
	re := a.get(c.Request.Method, c.FullPath())
	re.Name, re.Comment = getFuncComment()
	re.Req = ParmaItem(arg)
}
func (a *apiItems) setRes(c *gin.Context, arg any) {
	re := a.get(c.Request.Method, c.FullPath())
	re.Res = ParmaItem(arg)
}
func (a *apiItems) get(method, url string) *apiItem {
	for _, v := range *a {
		if v.Method == method && v.Url == url {
			return v
		}
	}
	return nil
}
func (a *apiItems) Add(method, url string) {
	*a = append(*a, &apiItem{
		Method: method,
		Url:    url,
	})
}
func (a *apiItems) cate() (res cateApiItems) {
	for _, item := range *a {
		cates := strings.Split(item.Url, "/")
		var re *cateApiItem
		for i, v := range cates[1:] {
			if i == 0 {
				re = res.init(v)
			} else {
				re = re.Children.init(v)
			}
			re.apiCount += 1
		}
		re.Apis = append(re.Apis, item)

	}
	sortCateApiItem(res)
	return
}
func (a *cateApiItems) init(name string) (re *cateApiItem) {
	for _, v := range *a {
		if v.Name == name {
			return v
		}
	}
	re = &cateApiItem{Name: name}
	*a = append(*a, re)
	return
}
func (a *cateApiItems) apis() (re apiItems) {
	for _, v := range *a {
		re = append(re, v.Apis...)
		re = append(re, v.Children.apis()...)
	}
	return
}
func sortCateApiItem(a cateApiItems) (res cateApiItems, apis apiItems) {
	if len(a) == 1 {
		a = a[0].Children
	}
	for _, v := range a {
		if v.apiCount == 1 && v.Children == nil {
			apis = append(apis, v.Apis...)
			continue
		} else if v.apiCount-len(v.Apis) <= 2 || v.apiCount-len(v.Children) <= 2 {
			v.Apis = append(v.Apis, v.Children.apis()...)
			v.Children = nil
		} else {
			var api apiItems
			v.Children, api = sortCateApiItem(v.Children)
			v.Apis = append(v.Apis, api...)
		}
		res = append(res, v)
	}
	return
}

type cateApiItem struct {
	Name     string
	Apis     apiItems
	Children cateApiItems
	apiCount int
}

type apiItem struct {
	Name    string
	Method  string
	Url     string
	Comment string
	Req     []structItem
	Res     []structItem
}

func (i apiItem) MD() string {
	re := "接口地址:" + i.Url + "  \n请求方式:" + i.Method + "  \n接口说明:" + i.Comment + "  \n请求参数:  \n  \n|参数名|必选|类型|说明|\n|--|--|--|--|\n"
	for _, v := range i.Req {
		re += v.MD(true) + "  \n"
	}
	re += "返回参数:  \n  \n|参数名|类型|说明|\n|--|--|--|\n"
	for _, v := range i.Res {
		re += v.MD(false) + "  \n"
	}
	return re
}
