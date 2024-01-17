package doc

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"strings"

	"gkk/expect"

	"github.com/gin-gonic/gin"
)

const apiMarkDown = "apis.md"

var docFlag bool
var apis apiItems

func Init(r *gin.Engine) {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		fmt.Fprintf(gin.DefaultWriter, "%-6s %-40s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
		apis.Add(httpMethod, absolutePath)
	}
}
func SetRes(c *gin.Context, data any) {
	if docFlag {
		apis.setRes(c, data)
	}
}
func SetReq(c *gin.Context, data any) {
	if docFlag {
		apis.setReq(c, data)
	}
}

func Handler(c *gin.Context) any {
	if c.Request.Method == "GET" {
		switch c.Query("out") {
		case "markdown":
			write2MD(c.Query("path"))
			c.File(apiMarkDown)
			return nil
		default:
			res := write2api(c.Query("path"))
			return res.cate()
		}
	} else {
		docFlag = !docFlag
		return apis.cate()
	}
}

func getFuncComment() (name, comments string) {
	pc, server, _, _ := runtime.Caller(3)
	name = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1]
	if src, err := parser.ParseFile(token.NewFileSet(), server, nil, parser.ParseComments); err == nil {
		for _, v := range src.Decls {
			if f, ok := v.(*ast.FuncDecl); ok && name == f.Name.Name {
				if f.Doc != nil {
					for i, comment := range f.Doc.List {
						switch i {
						case 0:
							name = strings.TrimLeft(comment.Text, "/")
						default:
							comments += strings.TrimLeft(comment.Text, "/")
						}
					}
				}
				return
			}
		}
	}
	return
}

func write2MD(url string) {
	file, err := os.OpenFile(apiMarkDown, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	expect.PEM(err, "创建失败")
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, item := range apis {
		if strings.Contains(item.Url, url) {
			write.WriteString(item.MD() + "  \n")
		}
	}
	write.Flush()
}

func write2api(url string) apiItems {
	res := apiItems{}
	for _, item := range apis {
		if strings.Contains(item.Url, url) {
			res = append(res, item)
		}
	}
	return res
}
