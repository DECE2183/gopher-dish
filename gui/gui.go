package gui

import (
	"fmt"
	"time"
	"unicode"

	"gopher-dish/gui/widgets"
	"gopher-dish/world"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

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

	btnPlay := widgets.NewButton("play", pixel.V(5, 29), pixel.V(60, 30))
	btnStop := widgets.NewButton("stop", pixel.V(70, 29), pixel.V(60, 30))
	btnNormal := widgets.NewButton("normal", pixel.V(150, 29), pixel.V(60, 30))
	btnEnergy := widgets.NewButton("energy", pixel.V(215, 29), pixel.V(60, 30))
	btnAge := widgets.NewButton("age", pixel.V(280, 29), pixel.V(60, 30))

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

		if btnPlay.Draw(win) {
			wd.world.TickInterval = 15 * time.Millisecond
		}
		if btnStop.Draw(win) {
			wd.world.TickInterval = 0
		}

		if btnNormal.Draw(win) {
			wd.Filter = W_FILTER_DISABLE
		}
		if btnEnergy.Draw(win) {
			wd.Filter = W_FILTER_ENERGY
		}
		if btnAge.Draw(win) {
			wd.Filter = W_FILTER_AGE
		}

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
	txt.Color = pixel.RGB(0.53, 0.53, 0.53)
	fmt.Fprintf(txt, "FPS: % 3d | Population: % 8d | Day: % 8d | Year: % 8d | Epoch: % 8d", world.Framerate, len(world.Objects), world.Ticks, world.Year, world.Epoch)
}
