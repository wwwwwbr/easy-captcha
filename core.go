package easyCaptcha

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

const (
	// 去除容易混淆的字母数字
	Seed_No_Similar_Alphanumeric = "abcdefghjkmnprstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ2345678"
	Seed_Alphanumeric            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Seed_Alphabetic              = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Seed_Numeric                 = "0123456789"
)

type Captcha interface {
	SetLineNum(int)
	SetBackgroundColor(color.Color)
	SetSeed(seed string)
	SetFont(captchaFont)

	Text() string
	GetBytes() ([]byte, error)
	SaveFile(filename string) error
	Base64() (string, error)
	generator()
}

type baseCaptcha struct {
	generated bool // 是否已生成
	width     int  // 宽度
	height    int  // 高度
	size      int  // 长度

	bgColor color.Color // 背景
	font    captchaFont // 字体
	lineNum int         // 干扰线数量
	text    string      // 对应的字符
	seed    string      // 生成字符的种子
}

func (c *baseCaptcha) SetLineNum(n int) {
	/*
		if n < 0 {
			n = 2
		}
		if n > 10 {
			n = 10
		}
	*/
	c.lineNum = n
}
func (c *baseCaptcha) SetBackgroundColor(col color.Color) {
	c.bgColor = col
}

func (c *baseCaptcha) SetSeed(seed string) {
	c.seed = seed
}
func (c *baseCaptcha) SetFont(f captchaFont) {
	c.font = f
}

// 绘制贝塞尔曲线
func drawBezierLine(dc *gg.Context, num int) {
	height := dc.Height()
	width := dc.Width()

	for i := 0; i < num; i++ {
		x1 := 5.0
		y1 := float64(randNumber(5, height/2))
		x2 := float64(width - 5)
		y2 := float64(randNumber(height/2, height-5))
		ctrlX := float64(randNumber(width/4, width/4*3))
		ctrlY := float64(randNumber(5, height-5))

		if r.Intn(2) == 0 {
			ty := y1
			y1 = y2
			y2 = ty
		}

		// 线条颜色
		dc.SetColor(randomColor())
		// 三阶贝塞尔曲线
		dc.CubicTo(x1, y1, ctrlX, ctrlY, x2, y2)

		dc.Stroke()
	}
}

// 绘制干扰泡泡
func drewCircle(dc *gg.Context, num int) {
	height := dc.Height()
	width := dc.Width()
	for range num {
		// 半径
		d := r.Intn(6) + 2
		// 随机位置
		x := r.Intn(width - d)
		y := r.Intn(height - d)
		dc.SetColor(color.Black) //黑色
		dc.DrawCircle(float64(x), float64(y), float64(d))
		dc.Stroke()
	}
}

func randomString(length int) string {
	return randomStringBySeed(length, Seed_No_Similar_Alphanumeric)
}

func randomStringBySeed(length int, seed string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = seed[r.Intn(len(seed))]
	}
	return string(b)
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randNumber(min, max int) int {
	return r.Intn(max-min) + min
}

func randomColorA(a uint8) color.RGBA {
	rc := r.Intn(256)
	g := r.Intn(256)
	b := r.Intn(256)
	return color.RGBA{uint8(rc), uint8(g), uint8(b), a}
}

func randomColor() color.RGBA {
	a := r.Intn(256)
	return randomColorA(uint8(a))
}
