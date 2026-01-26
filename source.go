package easyCaptcha

import (
	"embed"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//go:embed fonts/**
var SourceFS embed.FS

type captchaFont struct {
	path  string
	fType int // 0-内置 1-外部引用
}

var (
	ACTIONJ  = captchaFont{path: "actionj.ttf", fType: 0}
	EPI_LOG  = captchaFont{path: "epilog.ttf", fType: 0}
	FRESNEL  = captchaFont{path: "fresnel.ttf", fType: 0}
	HEADACHE = captchaFont{path: "headache.ttf", fType: 0}
	LEXO     = captchaFont{path: "lexo.ttf", fType: 0}
	PREFIX   = captchaFont{path: "prefix.ttf", fType: 0}
	PROGBOT  = captchaFont{path: "progbot.ttf", fType: 0}
	RANSOM   = captchaFont{path: "ransom.ttf", fType: 0}
	ROBOT    = captchaFont{path: "robot.ttf", fType: 0}
	SCANDAL  = captchaFont{path: "scandal.ttf", fType: 0}
)

func NewCaptchaFont(fontPath string) captchaFont {
	return captchaFont{path: fontPath, fType: 1}
}

func (c captchaFont) ParseFont() *truetype.Font {
	var bytes []byte
	var err error
	if c.fType == 1 {
		bytes, err = os.ReadFile(c.path)
	} else {
		bytes, err = SourceFS.ReadFile("fonts/" + string(c.path))
	}
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(bytes)
	if err != nil {
		panic(err)
	}
	return font
}

var HanZiLib = []rune("的是不了在人有我他这个们中国大上为和地到以说时要就出会可也用学生很好看语文字又都从自前后方向如果得而与本去来之发也些里思想情况感觉间新旧长短高低白红黄绿蓝黑少多每次定总般")
