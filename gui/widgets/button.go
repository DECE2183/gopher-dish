package widgets

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Button struct {
	bounds pixel.Rect
	label  string

	normalDrawer  *imdraw.IMDraw
	hoverDrawer   *imdraw.IMDraw
	pressedDrawer *imdraw.IMDraw
	labelDrawer   *text.Text
}

func NewButton(label string, pos, size pixel.Vec) *Button {
	b := &Button{
		label: label,
	}

	secDot := pos.Add(size)
	b.bounds = pixel.R(pos.X, pos.Y, secDot.X, secDot.Y)

	b.labelDrawer = text.New(b.bounds.Center(), fontAtlas)
	labelBounds := b.labelDrawer.BoundsOf(label)
	b.labelDrawer.Dot = b.labelDrawer.Orig.Sub(pixel.V(math.Round(labelBounds.W()/2), math.Round(labelBounds.H()/4)))
	b.labelDrawer.WriteString(label)

	b.normalDrawer = imdraw.New(nil)
	b.normalDrawer.Color = colorButton
	b.normalDrawer.Push(b.bounds.Min, b.bounds.Max)
	b.normalDrawer.Rectangle(0)

	b.hoverDrawer = imdraw.New(nil)
	b.hoverDrawer.Color = colorButtonHover
	b.hoverDrawer.Push(b.bounds.Min, b.bounds.Max)
	b.hoverDrawer.Rectangle(0)

	b.pressedDrawer = imdraw.New(nil)
	b.pressedDrawer.Color = colorButtonPressed
	b.pressedDrawer.Push(b.bounds.Min, b.bounds.Max)
	b.pressedDrawer.Rectangle(0)

	return b
}

func (b *Button) SetBounds(bounds pixel.Rect) {
	b.bounds = bounds
}

func (b *Button) SetPos(pos pixel.Vec) {
	b.bounds.Max = pos.Add(b.bounds.Size())
	b.bounds.Min = pos
}

func (b *Button) SetSize(size pixel.Vec) {
	b.bounds.Max = size.Add(b.bounds.Min)
}

func (b *Button) Draw(win *pixelgl.Window) bool {
	b.normalDrawer.Clear()

	defer func() {
		b.normalDrawer.Push(b.bounds.Min, b.bounds.Max)
		b.normalDrawer.Rectangle(0)
		b.normalDrawer.Draw(win)
		b.labelDrawer.Orig = b.bounds.Center()
		b.labelDrawer.Draw(win, pixel.IM)
	}()

	if b.bounds.Contains(win.MousePosition()) {
		if win.Pressed(pixelgl.MouseButtonLeft) {
			b.normalDrawer.Color = colorButtonPressed
			return true
		}
		b.normalDrawer.Color = colorButtonHover
		return false
	}

	b.normalDrawer.Color = colorButton
	return false
}
