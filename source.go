package easyCaptcha

import (
	"embed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//go:embed fonts/**
var SourceFS embed.FS

type captchaFont string

const (
	ACTIONJ  captchaFont = "actionj.ttf"
	EPI_LOG  captchaFont = "epilog.ttf"
	FRESNEL  captchaFont = "fresnel.ttf"
	HEADACHE captchaFont = "headache.ttf"
	LEXO     captchaFont = "lexo.ttf"
	PREFIX   captchaFont = "prefix.ttf"
	PROGBOT  captchaFont = "progbot.ttf"
	RANSOM   captchaFont = "ransom.ttf"
	ROBOT    captchaFont = "robot.ttf"
	SCANDAL  captchaFont = "scandal.ttf"
)

func (c captchaFont) GetBytes() []byte {
	bytes, err := SourceFS.ReadFile("fonts/" + string(c))
	if err != nil {
		panic(err)
	}
	return bytes
}

func (c captchaFont) ParseFont() *truetype.Font {
	font, err := freetype.ParseFont(c.GetBytes())
	if err != nil {
		panic(err)
	}
	return font
}
