package orc

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"github.com/otiai10/gosseract"
	"image"
	"image/color"
	"image/png"
)

func Recognize(data []byte) (res string, err error) {
	data, err = RemoveBackground(data)
	if err != nil {
		logs.Error(err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	err = client.SetImageFromBytes(data)
	if err != nil {
		logs.Error(err)
		return
	}
	res, err = client.Text()
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

//图片二值化处理
func RemoveBackground(data []byte) (res []byte, err error) {

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		logs.Error(err)
		return
	}

	bounds := img.Bounds()
	// 定义一个临界值 用于二值化界限
	var threshold uint8 = 150

	// 获取图片界限：image.Bounds
	// 获取某点像素 image.At  image.RGBAAt
	// 设置某点像素 image.Set image.SetRGBA

	imgSet := image.NewRGBA(bounds)

	for x := 1; x < bounds.Max.X; x++ {
		for y := 1; y < bounds.Max.Y; y++ {
			oldPixel := img.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			nr := uint8(r)
			ng := uint8(g)
			nb := uint8(b)
			if nr/3+ng/3+nb/3 <= threshold {
				pixel := color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}
				imgSet.Set(x, y, pixel)
			}
		}
	}

	buf := bytes.Buffer{}

	err = png.Encode(&buf, imgSet)
	if err != nil {
		logs.Error(err)
		return
	}
	res = buf.Bytes()

	return
}
