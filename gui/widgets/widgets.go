package widgets

import (
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

var (
	colorOutline       = pixel.RGB(0.4, 0.4, 0.4)
	colorButton        = pixel.RGB(0.3, 0.74, 1)
	colorButtonHover   = pixel.RGB(0.34, 0.84, 1)
	colorButtonPressed = pixel.RGB(0.25, 0.74, 1)
)

var (
	fontAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
)
