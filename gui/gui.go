package gui

import (
	"fmt"
	"time"
	"unicode"

	"gopher-dish/object"
	"gopher-dish/world"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type WorldDrawer struct {
	world         *world.World
	canvas        *pixelgl.Canvas
	baseDrawer    *imdraw.IMDraw
	objectsDrawer *imdraw.IMDraw
	matrix        pixel.Matrix
	bounds        pixel.Rect
	zoom          float64
}

var (
	updateTimerInterval time.Duration
	worldToDraw         *world.World
	textAtlas           *text.Atlas
)

const (
	DefaultZoomValue = 5
)

func Run(updateInterval time.Duration, world *world.World) {
	updateTimerInterval = updateInterval
	worldToDraw = world

	pixelgl.Run(initGUI)
}

func initGUI() {
	cfg := pixelgl.WindowConfig{
		Title:     "gopher-dish",
		Bounds:    pixel.R(0, 0, 1024, 480),
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	wd := NewWorldDrawer(worldToDraw)
	wd.IncZoom(DefaultZoomValue, pixel.V(0, 0))
	wd.Move(win.Bounds().Center())

	textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
	statusText := text.New(pixel.V(12, 8), textAtlas)

	statusPanelBg := imdraw.New(nil)
	createStatusBar(statusPanelBg, win.Bounds().Max.X, 24)

	var moveVec pixel.Vec
	var scrollVec pixel.Vec

	for !win.Closed() {
		if win.Pressed(pixelgl.MouseButtonRight) {
			if win.JustPressed(pixelgl.MouseButtonRight) {
				moveVec = win.MousePosition()
			} else {
				wd.Move(win.MousePosition().Sub(moveVec))
				moveVec = win.MousePosition()
			}
		}

		scrollVec = win.MouseScroll()
		if scrollVec.Y != 0 {
			wd.IncZoom(scrollVec.Y, win.Bounds().Center().Sub(win.MousePosition()))
		}

		win.Clear(colornames.White)
		wd.Draw(win)
		createStatusBar(statusPanelBg, win.Bounds().Max.X, 24)
		printStatus(statusText, wd.world, win.Bounds().Max.X, 24)
		statusPanelBg.Draw(win)
		statusText.Draw(win, pixel.IM)
		win.Update()
	}
}

func createStatusBar(imd *imdraw.IMDraw, width, height float64) {
	imd.Clear()
	imd.Color = pixel.RGB(0.83, 0.83, 0.83)
	imd.Push(pixel.V(0, height), pixel.V(width, height))
	imd.Line(1)
	imd.Color = pixel.RGB(0.92, 0.92, 0.92)
	imd.Push(pixel.V(0, 0), pixel.V(width, height-1))
	imd.Rectangle(0)
}

func printStatus(txt *text.Text, world *world.World, width, height float64) {
	if world.Ticks%10 != 0 {
		return
	}

	txt.Clear()
	txt.Color = pixel.RGB(0.73, 0.73, 0.73)
	fmt.Fprintf(txt, "Population: % 8d | Day: % 8d | Year: % 8d | Epoch: % 8d", len(world.Objects), world.Ticks, world.Year, world.Epoch)
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
	wd.bounds.Max = pixel.Vec{X: float64(wd.world.Width) * wd.zoom, Y: float64(wd.world.Height) * wd.zoom}
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
	sunlightColor := pixel.RGB(1, 0.83, 0.3)
	mineralsColor := pixel.RGB(0.3, 0.74, 1)

	wd.baseDrawer.Clear()

	for x := int32(0); x < int32(wd.world.Width); x++ {
		for y := int32(0); y < int32(wd.world.Height); y++ {
			sunlight := wd.world.GetSunlightAtPosition(object.Position{x, y})
			minerals := wd.world.GetMineralsAtPosition(object.Position{x, y})

			pixelColor := sunlightColor.Mul(pixel.Alpha(float64(sunlight) / 85))
			pixelColor = pixelColor.Add(mineralsColor.Mul(pixel.Alpha(float64(minerals) / 85)))

			wd.baseDrawer.Color = pixelColor

			var posX, posY = float64(x) * wd.zoom, wd.bounds.Max.Y - float64(y)*wd.zoom
			wd.baseDrawer.Push(pixel.V(posX, posY), pixel.V(posX+wd.zoom, posY-wd.zoom))
			wd.baseDrawer.Rectangle(0)
		}
	}
}

func (wd *WorldDrawer) DrawObjects() {
	// if <-wd.world.PlacesUpdated != true {
	// 	return
	// }

	wd.objectsDrawer.Clear()

	for x := 0; x < int(wd.world.Width); x++ {
		for y := 0; y < int(wd.world.Height); y++ {
			if wd.world.Places[x][y] == nil {
				continue
			}

			var posX, posY = float64(x) * wd.zoom, wd.bounds.Max.Y - float64(y)*wd.zoom
			wd.objectsDrawer.EndShape = imdraw.NoEndShape

			wd.objectsDrawer.Color = colornames.Black
			wd.objectsDrawer.Push(pixel.V(posX, posY), pixel.V(posX+wd.zoom, posY-wd.zoom))
			wd.objectsDrawer.Rectangle(1)

			wd.objectsDrawer.Color = colornames.Green
			wd.objectsDrawer.Push(pixel.V(posX, posY-1), pixel.V(posX+wd.zoom-1, posY-wd.zoom))
			wd.objectsDrawer.Rectangle(0)
		}
	}

	// wd.world.PlacesDrawn <- true
}

func (wd *WorldDrawer) Draw(t pixel.Target) {
	wd.canvas.Clear(colornames.White)
	wd.baseDrawer.Draw(wd.canvas)
	wd.DrawObjects()
	wd.objectsDrawer.Draw(wd.canvas)
	wd.canvas.Draw(t, wd.matrix)
}
