package widgets

import (
	"math"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

func Button(target pixel.Target, label string, pos pixel.Vec, size pixel.Vec, mousePos pixel.Vec, pressed bool) bool {
	secDot := pos.Add(size)

	drawer := imdraw.New(nil)
	drawer.EndShape = imdraw.RoundEndShape

	labelText := text.New(pos.Add(secDot).Scaled(0.5), fontAtlas)
	labelBounds := labelText.BoundsOf(label)
	labelText.Dot = labelText.Orig.Sub(pixel.V(math.Round(labelBounds.W()/2), math.Round(labelBounds.H()/4)))
	labelText.WriteString(label)

	if mousePos.X > pos.X && mousePos.Y > pos.Y && mousePos.X < secDot.X && mousePos.Y < secDot.Y {
		if pressed {
			drawer.Color = colorButtonPressed
		} else {
			drawer.Color = colorButtonHover
		}
	} else {
		drawer.Color = colorButton
		pressed = false
	}

	drawer.Push(pos, secDot)
	drawer.Rectangle(0)
	drawer.Color = colorOutline
	drawer.Push(pos, secDot)
	drawer.Rectangle(2)

	drawer.Draw(target)
	labelText.Draw(target, pixel.IM)

	return pressed
}
