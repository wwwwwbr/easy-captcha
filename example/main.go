package main

import (
	"fmt"

	"github.com/wwwwwbr/easy-captcha"
)

const GifPath = "./build/captcha.gif"
const PngPath = "./build/captcha.png"

func main() {
	captcha := easyCaptcha.NewGifCaptcha(120, 60, 4)
	err := captcha.SaveFile(GifPath)
	fmt.Println(err)

	text := captcha.Text()
	fmt.Println(text)

	base64, err := captcha.Base64()
	fmt.Println(err)
	fmt.Println(base64)

	fmt.Println("======================")

	simpleCaptcha := easyCaptcha.NewSimpleCaptcha(120, 60, 4)
	//simpleCaptcha.SetLineNum(4)
	simpleCaptcha.SetSeed("qwertyuiop123")
	//simpleCaptcha.SetBackgroundColor(color.RGBA{100, 100, 100, 155})
	//simpleCaptcha.SetFont(easyCaptcha.LEXO)
	err = simpleCaptcha.SaveFile(PngPath)
	fmt.Println(err)
	text = simpleCaptcha.Text()
	fmt.Println(text)
	base64, err = simpleCaptcha.Base64()
	fmt.Println(err)
	fmt.Println(base64)

}
