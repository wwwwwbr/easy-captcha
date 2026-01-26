package main

import (
	"fmt"

	"github.com/wwwwwbr/easy-captcha"
)

const GifPath = "./build/captcha.gif"
const PngPath = "./build/captcha.png"
const HzPath = "./build/captcha_hz.png"
const CalcPath = "./build/captcha_calc.png"

func main() {
	exampleChinese()
	exampleSimple()
	exampleGif()
	exampleCalc()
}

func exampleCalc() {
	calculationCaptcha := easyCaptcha.NewCalculationCaptcha(240, 100, 3)
	text := calculationCaptcha.Text()
	fmt.Println("result:", text)
	err := calculationCaptcha.SaveFile(CalcPath)
	if err != nil {
		fmt.Println(err)
	}
}

func exampleChinese() {

	fp := "./example/SIMKAI.TTF"

	captcha := easyCaptcha.NewChineseCaptcha(200, 100, 4, fp)
	captcha.SetSeed("零壹贰叁肆伍陆柒捌玖拾佰仟万亿")

	text := captcha.Text()
	fmt.Println(text)

	err := captcha.SaveFile(HzPath)
	if err != nil {
		fmt.Println(err)
	}
}

func exampleSimple() {
	simpleCaptcha := easyCaptcha.NewSimpleCaptcha(120, 60, 4)
	//simpleCaptcha.SetLineNum(4)
	simpleCaptcha.SetSeed("qwertyuiop123")
	//simpleCaptcha.SetBackgroundColor(color.RGBA{100, 100, 100, 155})
	//simpleCaptcha.SetFont(easyCaptcha.LEXO)
	err := simpleCaptcha.SaveFile(PngPath)
	fmt.Println(err)
	text := simpleCaptcha.Text()
	fmt.Println(text)
	base64, err := simpleCaptcha.Base64()
	fmt.Println(err)
	fmt.Println(base64)
}

func exampleGif() {
	gifCaptcha := easyCaptcha.NewGifCaptcha(120, 60, 4)
	err := gifCaptcha.SaveFile(GifPath)
	fmt.Println(err)

	text := gifCaptcha.Text()
	fmt.Println(text)

	base64, err := gifCaptcha.Base64()
	fmt.Println(err)
	fmt.Println(base64)
}
