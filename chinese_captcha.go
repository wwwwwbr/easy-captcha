package easyCaptcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type chineseCaptcha struct {
	hanziArr []string
	baseCaptcha
	img *image.Image
}

func NewChineseCaptcha(width, height, size int, fontPath string) Captcha {
	sc := &chineseCaptcha{
		baseCaptcha: baseCaptcha{
			bgColor:   color.White,
			width:     width,
			height:    height,
			size:      size,
			font:      NewCaptchaFont(fontPath),
			lineNum:   2,
			generated: false,
		},
	}
	var c Captcha = sc
	return c
}

func (c *chineseCaptcha) randomHanzi() {
	c.hanziArr = make([]string, c.size)
	var hc int
	var hzLib []rune
	// 使用自定的字符串
	if c.seed != "" {
		hzLib = []rune(c.seed)
	} else {
		hzLib = HanZiLib
	}

	hc = len(hzLib)
	for i := range c.hanziArr {
		c.hanziArr[i] = string(hzLib[r.Intn(hc)])
	}
	c.text = strings.Join(c.hanziArr, "")
}

func (c *chineseCaptcha) generator() {
	c.randomHanzi()
	// 设置背景颜色
	bg := gg.NewContext(c.width, c.height)
	bg.SetColor(c.bgColor)
	bg.Clear()

	// 绘制随机贝塞尔线
	drawBezierLine(bg, c.lineNum)

	// 绘制文字 + 干扰泡泡
	width := c.width
	num := len(c.hanziArr)
	minWidth := width / num
	f := c.font.ParseFont()
	fontHeight := float64(c.height)
	face := truetype.NewFace(f, &truetype.Options{
		Size: fontHeight / 2,
	})

	// 字体颜色
	colors := make([]color.RGBA, num)
	for i := 0; i < num; i++ {
		colors[i] = randomColorA(255)
	}

	// 依次写入文字
	dc := gg.NewContext(bg.Width(), bg.Height())
	// bg
	dc.DrawImage(bg.Image(), 0, 0)

	// 绘制随机小圆圈
	drewCircle(dc, 2)

	// 文字
	dc.SetFontFace(face)

	// 兼容汉字
	for i := 0; i < num; i++ {
		ss := c.hanziArr[i]
		dc.SetColor(colors[i])
		dc.DrawStringAnchored(ss, float64(i*minWidth)+20, fontHeight/2, 0.5, 0.5)
	}
	// 生成图片
	img := dc.Image()
	c.img = &img

	c.generated = true
}

func (c *chineseCaptcha) Text() string {
	c.randomHanzi()
	if !c.generated {
		c.generator()
	}
	return c.text
}

func (c *chineseCaptcha) GetBytes() ([]byte, error) {
	if !c.generated {
		c.generator()
	}

	var buf bytes.Buffer
	// 编码
	err := png.Encode(&buf, *c.img)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *chineseCaptcha) SaveFile(savePath string) error {
	if !c.generated {
		c.generator()
	}
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	err = png.Encode(file, *c.img)
	return err
}

func (c *chineseCaptcha) Base64() (string, error) {
	if !c.generated {
		c.generator()
	}
	imgBytes, err := c.GetBytes()
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBytes), nil
}
