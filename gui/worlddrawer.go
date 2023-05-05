package gui

import (
	"gopher-dish/object"
	"gopher-dish/world"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	colorOutline  = pixel.RGB(0.1, 0.1, 0.1)
	colorSunlight = pixel.RGB(1, 0.83, 0.3)
	colorMinerals = pixel.RGB(0.3, 0.74, 1)

	colorObjectLively = pixel.RGB(0.38, 0.58, 0.27)

	colorObjectEnergyMax = pixel.RGB(1, 0.91, 0)
	colorObjectEnergyMin = pixel.RGB(1, 0, 0)
	colorObjectAgeMax    = pixel.RGB(0, 0.015, 0.52)
	colorObjectAgeMin    = pixel.RGB(0.1, 0.83, 0.78)

	colorObjectPickable = pixel.RGB(0.41, 0.074, 0.25)
	colorObjectMovalble = pixel.RGB(0.4, 0.3, 0.47)
	colorObjectUnknown  = colornames.Magenta
)

type WorldFilter byte

const (
	W_FILTER_DISABLE    WorldFilter = iota
	W_FILTER_ENERGY     WorldFilter = iota
	W_FILTER_AGE        WorldFilter = iota
	W_FILTER_FOOD_TYPE  WorldFilter = iota
	W_FILTER_GENERATION WorldFilter = iota
)

type WorldDrawer struct {
	Filter WorldFilter

	world         *world.World
	canvas        *pixelgl.Canvas
	baseDrawer    *imdraw.IMDraw
	objectsDrawer *imdraw.IMDraw
	matrix        pixel.Matrix
	bounds        pixel.Rect
	zoom          float64
}

func NewWorldDrawer(world *world.World) *WorldDrawer {
	wd := &WorldDrawer{
		world:         world,
		baseDrawer:    imdraw.New(nil),
		objectsDrawer: imdraw.New(nil),
		canvas:        pixelgl.NewCanvas(pixel.R(0, 0, float64(world.Width)*DefaultZoomValue, float64(world.Height)*DefaultZoomValue)),
		matrix:        pixel.IM,
	}
	return wd
}

func (wd *WorldDrawer) Move(p pixel.Vec) {
	wd.matrix = wd.matrix.Moved(p)
}

func (wd *WorldDrawer) IncZoom(level float64, vec pixel.Vec) {
	wd.zoom += level
	if wd.zoom < 1 {
		wd.zoom = 1
	} else if wd.zoom > 32 {
		wd.zoom = 32
	}

	oldBoundsMax := wd.bounds.Max
	wd.bounds.Max = pixel.Vec{X: float64(wd.world.Width-1) * wd.zoom, Y: float64(wd.world.Height-1) * wd.zoom}
	wd.canvas.SetBounds(wd.bounds)
	wd.DrawBase()

	if vec.X != 0 && vec.Y != 0 {
		if level > 0 {
			wd.Move(vec.ScaledXY(pixel.V(wd.bounds.Max.X/oldBoundsMax.X, wd.bounds.Max.Y/oldBoundsMax.Y)).Scaled(0.25))
		} else {
			wd.Move(vec.ScaledXY(pixel.V(-oldBoundsMax.X/wd.bounds.Max.X, -oldBoundsMax.Y/wd.bounds.Max.Y)).Scaled(0.25))
		}
	}
}

func (wd *WorldDrawer) DrawBase() {
	wd.baseDrawer.Clear()

	for x := int32(0); x < int32(wd.world.Width); x++ {
		for y := int32(0); y < int32(wd.world.Height); y++ {
			sunlight := wd.world.GetSunlightAtPosition(object.Position{X: x, Y: y})
			minerals := wd.world.GetMineralsAtPosition(object.Position{X: x, Y: y})

			pixelColor := colorSunlight.Mul(pixel.Alpha(float64(sunlight) / 85))
			pixelColor = pixelColor.Add(colorMinerals.Mul(pixel.Alpha(float64(minerals) / 85)))

			wd.baseDrawer.Color = pixelColor

			var posX, posY = float64(x) * wd.zoom, wd.bounds.Max.Y - float64(y)*wd.zoom
			wd.baseDrawer.Push(pixel.V(posX, posY), pixel.V(posX+wd.zoom, posY-wd.zoom))
			wd.baseDrawer.Rectangle(0)
		}
	}

	wd.baseDrawer.Color = colorOutline
	wd.baseDrawer.Push(pixel.V(1, 1), wd.canvas.Bounds().Max)
	wd.baseDrawer.Rectangle(1)
}

func (wd *WorldDrawer) ComputeObjectColor(obj object.Object) color.Color {
	var _obj interface{} = obj
	switch o := _obj.(type) {
	case object.Lively:
		switch wd.Filter {
		case W_FILTER_DISABLE:
			return colorObjectLively
		case W_FILTER_ENERGY:
			return lerpColor(colorObjectEnergyMin, colorObjectEnergyMax, float64(o.GetEnergy())/128)
		case W_FILTER_AGE:
			return lerpColor(colorObjectAgeMin, colorObjectAgeMax, float64(o.GetAge())/128)
		}
	case object.Pickable:
		return colorObjectPickable
	case object.Movable:
		return colorObjectMovalble
	}

	return colorObjectUnknown
}

func (wd *WorldDrawer) DrawObjects() {
	wd.world.PlacesDrawMux.Lock()
	defer wd.world.PlacesDrawMux.Unlock()

	wd.objectsDrawer.Clear()

	for x := 0; x < int(wd.world.Width); x++ {
		for y := 0; y < int(wd.world.Height); y++ {
			if wd.world.Places[x][y] == nil {
				continue
			}

			var posX, posY = float64(x) * wd.zoom, wd.bounds.Max.Y - float64(y)*wd.zoom
			wd.objectsDrawer.EndShape = imdraw.NoEndShape

			if wd.zoom > 9 {
				// draw outline
				wd.objectsDrawer.Color = colorOutline
				wd.objectsDrawer.Push(pixel.V(posX, posY), pixel.V(posX+wd.zoom, posY-wd.zoom))
				wd.objectsDrawer.Rectangle(1)
				// draw cell
				wd.objectsDrawer.Color = wd.ComputeObjectColor(wd.world.Places[x][y])
				wd.objectsDrawer.Push(pixel.V(posX, posY-1), pixel.V(posX+wd.zoom-1, posY-wd.zoom))
				wd.objectsDrawer.Rectangle(0)
			} else {
				// if cells are too small, draw only cells
				wd.objectsDrawer.Color = wd.ComputeObjectColor(wd.world.Places[x][y])
				wd.objectsDrawer.Push(pixel.V(posX, posY), pixel.V(posX+wd.zoom, posY-wd.zoom))
				wd.objectsDrawer.Rectangle(0)
			}
		}
	}
}

func (wd *WorldDrawer) Draw(t pixel.Target) {
	wd.canvas.Clear(colornames.White)
	wd.baseDrawer.Draw(wd.canvas)
	wd.DrawObjects()
	wd.objectsDrawer.Draw(wd.canvas)
	wd.canvas.Draw(t, wd.matrix)
}

func lerpColor(a pixel.RGBA, b pixel.RGBA, v float64) pixel.RGBA {
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}

	return pixel.RGBA{
		R: a.R + (b.R-a.R)*v,
		G: a.G + (b.G-a.G)*v,
		B: a.B + (b.B-a.B)*v,
		A: 1.0,
	}
}