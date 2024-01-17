package gkk

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"strings"

	"gkk/api"
	"gkk/config"
	"gkk/expect"
	"gkk/req"
	"gkk/tool"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const DEFAULT_MAX_WIDTH float64 = 320
const DEFAULT_MAX_HEIGHT float64 = 240

// 计算图片缩放后的尺寸
func calculateRatioFit(srcWidth, srcHeight int) (int, int) {
	ratio := math.Min(DEFAULT_MAX_WIDTH/float64(srcWidth), DEFAULT_MAX_HEIGHT/float64(srcHeight))
	return int(math.Ceil(float64(srcWidth) * ratio)), int(math.Ceil(float64(srcHeight) * ratio))
}

func CreateThumb(args []byte) []byte {
	img, _, _ := image.Decode(bytes.NewBuffer(args))
	b := img.Bounds()
	w, h := calculateRatioFit(b.Max.X, b.Max.Y)
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	res := new(bytes.Buffer)
	png.Encode(res, m)
	return res.Bytes()
}

func Gin2Byte(c *gin.Context) (res []byte, name string) {
	file, err := c.FormFile("file")
	expect.PEM(err, "文件上传失败")
	name = file.Filename
	f, e := file.Open()
	expect.PEM(e, "文件打开失败")
	defer f.Close()
	res, err = ioutil.ReadAll(f)
	expect.PEM(err, "读取失败")
	return
}

func UploadThumb4QNUrl(url string) string {
	res := strings.Split(strings.TrimLeft(strings.TrimLeft(url, "https://"), "http://"), "/")
	if ll := len(res); ll > 2 {
		names := strings.Split(res[ll-1], ".")
		names[0] += "_thumb"
		res[ll-1] = strings.Join(names, ".")
		body := HttpGet(url)
		return UploadQNBase(tool.Base64(CreateThumb(body)), strings.Join(res[1:], "/"))
	}
	return ""
}

// 将网络图片地址上传到七牛
func UploadQNUrl(url, prefixPath string, names ...string) (re string) {
	defer func() {
		recover()
	}()
	body := HttpGet(url)
	pathName := ""
	if len(names) > 0 {
		pathName = names[0]
	}
	return UploadQN(tool.Base64(body), pathName, prefixPath)
}

// 将gin的文件保存到七牛
func UploadQNGin(prefixPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s req.NameValue
		api.BindParam(c, &s)
		file, _ := Gin2Byte(c)
		file = tool.Base64(file)
		api.RD(c, UploadQN(file, s.Name, prefixPath))
	}
}

func UploadQN(base64 []byte, pathName, prefixPath string) (re string) {
	if pathName == "" {
		// 没有名字
		pathName = tool.UUID()
	} else if !strings.Contains(pathName, "/") {
		// 有名字没目录
		pathName = "other/" + pathName
	} else {
		pathName = strings.TrimLeft(pathName, "/")
	}
	pathName = prefixPath + pathName
	return UploadQNBase(base64, pathName)
}

// base64的byte流上传到七牛
func UploadQNBase(base64 []byte, pathName string) (re string) {
	qn := config.Get().Qn
	putPolicy := storage.PutPolicy{Scope: qn.SC}
	zone, _ := storage.GetRegionByID(storage.RegionID(qn.ST))
	base64Uploader := storage.NewBase64Uploader(&storage.Config{
		Zone:          &zone,
		UseHTTPS:      false,
		UseCdnDomains: true,
	})
	err := base64Uploader.Put(
		context.Background(),
		&storage.PutRet{},
		putPolicy.UploadToken(qbox.NewMac(qn.AK, qn.SK)),
		pathName,
		base64,
		nil)
	if err != nil {
		Log.Error("图片上传失败", err.Error())
		re = ""
	} else {
		re = qn.Host + pathName
	}
	return
}
