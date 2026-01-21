package easyCaptcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type SimpleCaptcha struct {
	baseCaptcha
	img *image.Image
}

func NewSimpleCaptcha(width, height, size int) Captcha {
	sc := &SimpleCaptcha{
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
	var c Captcha = sc
	return c
}

// 生成图片
func (s *SimpleCaptcha) generator() {

	// 生成文字
	if s.seed != "" {
		s.text = randomStringBySeed(s.size, s.seed)
	} else {
		s.text = randomString(s.size)
	}

	// 设置背景颜色
	bg := gg.NewContext(s.width, s.height)
	bg.SetColor(s.bgColor)
	bg.Clear()

	// 绘制随机贝塞尔线
	drawBezierLine(bg, s.lineNum)

	// 绘制文字 + 干扰泡泡
	width := s.width
	num := len(s.text)
	minWidth := width / num
	f := s.font.ParseFont()
	fontHeight := float64(s.height)
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
	for i := 0; i < num; i++ {
		ss := s.text[i : i+1]
		dc.SetColor(colors[i])
		dc.DrawStringAnchored(ss, float64(i*minWidth)+20, fontHeight/2, 0.5, 0.5)
	}

	// 生成图片
	img := dc.Image()
	s.img = &img

	s.generated = true
}

func (s *SimpleCaptcha) Text() string {
	if !s.generated {
		s.generator()
	}
	return s.text
}

func (s *SimpleCaptcha) GetBytes() ([]byte, error) {
	if !s.generated {
		s.generator()
	}

	var buf bytes.Buffer
	// 编码
	err := png.Encode(&buf, *s.img)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *SimpleCaptcha) SaveFile(savePath string) error {
	if !s.generated {
		s.generator()
	}
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	err = png.Encode(file, *s.img)
	return err
}

func (s *SimpleCaptcha) Base64() (string, error) {
	if !s.generated {
		s.generator()
	}
	imgBytes, err := s.GetBytes()
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBytes), nil
}
