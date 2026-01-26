package easyCaptcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type calculationCaptcha struct {
	baseCaptcha
	img        *image.Image
	expression string
}

func NewCalculationCaptcha(width, height, size int) Captcha {
	sc := &calculationCaptcha{
		baseCaptcha: baseCaptcha{
			bgColor:   color.White,
			width:     width,
			height:    height,
			font:      ACTIONJ,
			size:      size,
			lineNum:   2,
			generated: false,
		},
	}
	var c Captcha = sc
	return c
}

var calc_op = []string{"+", "-", "*"}

// 随机表达式
func (c *calculationCaptcha) randomExpr() {
	if c.expression != "" {
		return
	}
	numRange, numOp := 99, 2
	if c.size > 2 {
		numRange, numOp = 9, 3
	}
	for i := range c.size {
		n := r.Intn(numRange) + 1
		if i != 0 {
			// 追加符号
			op := calc_op[r.Intn(numOp)]
			c.expression += op
		}
		c.expression += strconv.Itoa(n)
	}
	result, _ := parseAndCalculate(c.expression)
	// 计算结果
	c.text = strconv.Itoa(int(result))
	// 乘号替换
	c.expression = strings.ReplaceAll(c.expression, "*", "x")
	// 拼接尾缀
	c.expression += "=?"
}

func (c *calculationCaptcha) generator() {
	// 生成表达式
	c.randomExpr()

	// 设置背景颜色
	bg := gg.NewContext(c.width, c.height)
	bg.SetColor(c.bgColor)
	bg.Clear()

	// 绘制随机贝塞尔线
	drawBezierLine(bg, c.lineNum)

	// 绘制文字 + 干扰泡泡
	width := c.width
	// 使用表达式绘制
	num := len(c.expression)
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
	for i := 0; i < num; i++ {
		ss := c.expression[i : i+1]
		dc.SetColor(colors[i])
		dc.DrawStringAnchored(ss, float64(i*minWidth)+20, fontHeight/2, 0.5, 0.5)
	}

	// 生成图片
	img := dc.Image()
	c.img = &img

	c.generated = true
}

func (c *calculationCaptcha) Text() string {
	c.randomExpr()
	if !c.generated {
		c.generator()
	}
	return c.text
}

func (c *calculationCaptcha) GetBytes() ([]byte, error) {
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

func (c *calculationCaptcha) SaveFile(savePath string) error {
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

func (c *calculationCaptcha) Base64() (string, error) {
	if !c.generated {
		c.generator()
	}
	imgBytes, err := c.GetBytes()
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBytes), nil
}

func calculateOperation(num1Str, op, num2Str string) (float64, error) {
	num1, err := strconv.ParseFloat(num1Str, 64)
	if err != nil {
		return 0, fmt.Errorf("无效数字 '%s': %w", num1Str, err)
	}
	num2, err := strconv.ParseFloat(num2Str, 64)
	if err != nil {
		return 0, fmt.Errorf("无效数字 '%s': %w", num2Str, err)
	}

	var res float64
	switch op {
	case "*":
		res = num1 * num2
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		res = num1 / num2
	case "+":
		res = num1 + num2
	case "-":
		res = num1 - num2
	default:
		return 0, fmt.Errorf("不支持的运算符: %s", op)
	}
	return res, nil
}

func evaluateSimpleExpression(expr string) (float64, error) {
	// 正则表达式用于提取 token:
	// `([+\-])` 匹配运算符 '+' 或 '-'
	// `(\-?\d+(\.\d+)?)` 匹配数字 (例如 "123", "-45", "3.14", "-0.5")
	tokenRe := regexp.MustCompile(`([+\-])|(\-?\d+(\.\d+)?)`)
	parts := tokenRe.FindAllString(expr, -1)

	if len(parts) == 0 {
		return 0, fmt.Errorf("无效或空表达式: %s", expr)
	}

	// 第一个 token 必须是数字 (可能带负号)
	result, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("无效的起始数字 '%s': %w", parts[0], err)
	}

	// 遍历剩余的 token，期望是 "运算符-数字" 对
	for i := 1; i < len(parts); i += 2 {
		op := parts[i]
		if i+1 >= len(parts) {
			return 0, fmt.Errorf("表达式格式错误: 运算符 '%s' 缺少右操作数", op)
		}
		numStr := parts[i+1]

		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("运算符 '%s' 后的数字无效 '%s': %w", op, numStr, err)
		}

		switch op {
		case "+":
			result += num
		case "-":
			result -= num
		default:
			return 0, fmt.Errorf("最终计算中遇到不支持的运算符: %s", op)
		}
	}
	return result, nil
}

// parseAndCalculate 解析并计算算术表达式
// 支持 +,-,*,/ 运算符和优先级，但不支持括号。
// 示例: "8-5*6", "10+2*3-4/2", "5*-2"
func parseAndCalculate(expr string) (float64, error) {
	// 移除所有空格
	currentExpr := strings.ReplaceAll(expr, " ", "")

	// 步骤 1: 优先处理乘法和除法运算
	// 正则表达式 `(\-?\d+(\.\d+)?)\s*([*/])\s*(\-?\d+(\.\d+)?)` 匹配:
	// - 第一个数字 (可能带负号和小数点)
	// - 运算符 '*' 或 '/'
	// - 第二个数字 (可能带负号和小数点)
	highPrecedenceRe := regexp.MustCompile(`(\-?\d+(\.\d+)?)\s*([*/])\s*(\-?\d+(\.\d+)?)`)

	for highPrecedenceRe.MatchString(currentExpr) {
		match := highPrecedenceRe.FindStringSubmatch(currentExpr)
		if len(match) < 6 {
			return 0, fmt.Errorf("高优先级运算解析错误: %s", match)
		}

		num1Str := match[1] // 第一个捕获组: 第一个数字
		op := match[3]      // 第三个捕获组: 运算符
		num2Str := match[4] // 第四个捕获组: 第二个数字

		res, err := calculateOperation(num1Str, op, num2Str)
		if err != nil {
			return 0, err
		}

		// 找到第一个匹配操作的起始和结束索引，然后替换它
		idx := highPrecedenceRe.FindStringIndex(currentExpr)
		if len(idx) == 2 {
			currentExpr = currentExpr[:idx[0]] + strconv.FormatFloat(res, 'f', -1, 64) + currentExpr[idx[1]:]
		} else {
			return 0, fmt.Errorf("内部错误: 无法找到正则表达式匹配索引")
		}
	}

	// 步骤 2: 处理剩余的加法和减法运算
	// 此时表达式应只包含数字、加号和减号。
	finalResult, err := evaluateSimpleExpression(currentExpr)
	if err != nil {
		return 0, fmt.Errorf("计算最终加减法时出错: %w", err)
	}

	return finalResult, nil
}
