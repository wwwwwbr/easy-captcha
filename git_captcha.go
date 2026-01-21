package easyCaptcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type GifCaptcha struct {
	baseCaptcha
	gifImg *gif.GIF
}

func NewGifCaptcha(width, height, size int) Captcha {
	gc := &GifCaptcha{
		baseCaptcha: baseCaptcha{
			bgColor:   color.White,
			width:     width,
			height:    height,
			size:      size,
			font:      ACTIONJ,
			lineNum:   2,
			generated: false,
		},
	}
	var c Captcha = gc
	return c
}

// 生成图片
func (g *GifCaptcha) generator() {

	// 生成文字
	if g.seed != "" {
		g.text = randomStringBySeed(g.size, g.seed)
	} else {
		g.text = randomString(g.size)
	}

	// 设置背景颜色
	bg := gg.NewContext(g.width, g.height)
	bg.SetColor(g.bgColor)
	bg.Clear()

	// 绘制随机贝塞尔线
	drawBezierLine(bg, g.lineNum)

	// 绘制文字 + 干扰泡泡
	width := g.width
	num := len(g.text)
	minWidth := width / num
	f := g.font.ParseFont()
	fontHeight := float64(g.height)
	face := truetype.NewFace(f, &truetype.Options{
		Size: fontHeight / 2,
	})

	// 字体颜色
	colors := make([]color.RGBA, num)
	for i := 0; i < num; i++ {
		colors[i] = randomColorA(255)
	}
	// gif帧
	images := make([]*image.Paletted, num)
	delays := make([]int, num)

	// 依次写入文字
	for c := 0; c < num; c++ {
		dc := gg.NewContext(bg.Width(), bg.Height())
		// bg
		dc.DrawImage(bg.Image(), 0, 0)

		// 绘制随机小圆圈
		drewCircle(dc, 2)

		// 文字
		dc.SetFontFace(face)
		for i := 0; i < num; i++ {
			if c == i {
				continue
			}
			ss := g.text[i : i+1]
			dc.SetColor(colors[i])
			dc.DrawStringAnchored(ss, float64(i*minWidth)+20, fontHeight/2, 0.5, 0.5)
		}

		// 绘制帧
		bounds := dc.Image().Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, bg.Image().Bounds(), dc.Image(), bounds.Min, draw.Src)

		images[c] = palettedImage
		delays[c] = 10
	}

	// 创建GIF
	g.gifImg = &gif.GIF{
		Image:     images, // make([]*image.Paletted, 0)
		Delay:     delays,
		LoopCount: 0,
	}

	g.generated = true
}

func (g *GifCaptcha) Text() string {
	if !g.generated {
		g.generator()
	}
	return g.text
}

func (g *GifCaptcha) GetBytes() ([]byte, error) {
	if !g.generated {
		g.generator()
	}

	var buf bytes.Buffer
	// 编码 GIF
	err := gif.EncodeAll(&buf, g.gifImg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GifCaptcha) SaveFile(savePath string) error {
	if !g.generated {
		g.generator()
	}
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	err = gif.EncodeAll(file, g.gifImg)
	return err
}

func (g *GifCaptcha) Base64() (string, error) {
	if !g.generated {
		g.generator()
	}
	imgBytes, err := g.GetBytes()
	if err != nil {
		return "", err
	}
	return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(imgBytes), nil
}
