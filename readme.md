# easy-captcha

## example
![png](example/captcha1.png)
![png](example/captcha2.png)
![gif](example/captcha.gif)
![gif](example/captcha_hz.png)
![gif](example/captcha_calc.png)

## Install

To use easy-captcha in your Go project, you can import it using the following command:

```shell
 go get github.com/wwwwwbr/easy-captcha
```

## Usage

To use the timezone conversion functions, first import this package:

```go
import (
    easyCaptcha "github.com/wwwwwbr/easy-captcha"
)
```
Then, call the function
```go
captcha := easyCaptcha.NewGifCaptcha(120, 60, 4) // gif
//captcha := easyCaptcha.NewSimpleCaptcha(120, 60, 4) // simple
//captcha := easyCaptcha.exampleChinese(120, 60, 4) // chinese
//captcha := easyCaptcha.NewCalculationCaptcha(120, 60, 3) // calculate

_ = captcha.SaveFile(GifPath)

text := captcha.Text()
fmt.Println(text)

base64, _ := captcha.Base64()
fmt.Println(base64)
	
```