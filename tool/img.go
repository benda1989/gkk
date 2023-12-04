package tool

import (
	"bytes"
	"gkk/expect"
	"image"
	"image/draw"
	"image/png"
	"io"
	"math"

	"github.com/nfnt/resize"
)

const DEFAULT_MAX_WIDTH int = 320
const DEFAULT_MAX_HEIGHT int = 240

func CreateThumb(r io.Reader) *bytes.Buffer {
	return CreateThumbSize(r, DEFAULT_MAX_WIDTH, DEFAULT_MAX_HEIGHT)
}

func CreateThumbSize(r io.Reader, width, height int) *bytes.Buffer {
	img, _, _ := image.Decode(r)
	b := img.Bounds()
	ratio := math.Min(float64(width)/float64(b.Max.X), float64(height)/float64(b.Max.Y))
	w, h := int(math.Ceil(float64(b.Max.X)*ratio)), int(math.Ceil(float64(b.Max.Y)*ratio))
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	res := new(bytes.Buffer)
	png.Encode(res, m)
	return res
}

func CombineImgs(imgs []*bytes.Buffer) *bytes.Buffer {
	sideLen := IntSqrt(len(imgs))
	if sideLen < 2 {
		sideLen = 2
	}
	var w, h int
	var canvas *image.RGBA
	for i, img := range imgs {
		imgt, _, err := image.Decode(img)
		if err != nil {
			return nil
		}
		width := imgt.Bounds().Max.X

		height := imgt.Bounds().Max.Y
		w1 := i % sideLen
		h1 := i / sideLen
		if i == 0 {
			if len(imgs) < 3 {
				canvas = image.NewRGBA(image.Rect(0, 0, width*2, height))
			} else {
				canvas = image.NewRGBA(image.Rect(0, 0, width*sideLen, height*sideLen))
			}
			w = width
			h = height
		} else {
			expect.PBM(w != width || h != height, "图片尺寸不一致")
		}

		draw.Draw(canvas, image.Rect(w1*width, h1*height, (w1+1)*width, (h1+1)*height), imgt, image.Point{0, 0}, draw.Src)
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, canvas)
	//outputFile, err := os.Create("output.png")
	//if err != nil {
	//	panic(err)
	//}
	//defer outputFile.Close()
	//err = png.Encode(outputFile, canvas)
	return buf
}

func IntSqrt(x int) int {
	// 使用牛顿迭代法求平方根
	if x == 0 {
		return 0
	}
	var sqrt = float64(x)
	for i := 0; i < 100; i++ {
		sqrt = (sqrt + float64(x)/sqrt) / 2
	}
	return int(sqrt)
}
