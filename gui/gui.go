package gui

import (
	"fmt"
	"os"
	"time"
	"unicode"

	"gopher-dish/gui/filepicker"
	"gopher-dish/gui/widgets"
	"gopher-dish/world"
	"gopher-dish/world/worldsaver"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	ticker      *time.Ticker
	worldToDraw *world.World
	textAtlas   *text.Atlas
)

func Run(updateInterval time.Duration, world *world.World) {
	ticker = time.NewTicker(updateInterval)
	worldToDraw = world
	pixelgl.Run(initGUI)
}

func initGUI() {
	cfg := pixelgl.WindowConfig{
		Title:     "gopher-dish",
		Bounds:    pixel.R(0, 0, 1736, 920),
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	wd := NewWorldDrawer(worldToDraw)
	wd.Move(win.Bounds().Center())

	btnPlay := widgets.NewButton("play", pixel.V(5, 29), pixel.V(60, 30))
	btnStop := widgets.NewButton("stop", pixel.V(70, 29), pixel.V(60, 30))

	btnNormal := widgets.NewButton("normal", pixel.V(150, 29), pixel.V(60, 30))
	btnHealth := widgets.NewButton("health", pixel.V(215, 29), pixel.V(60, 30))
	btnEnergy := widgets.NewButton("energy", pixel.V(280, 29), pixel.V(60, 30))
	btnAge := widgets.NewButton("age", pixel.V(345, 29), pixel.V(60, 30))

	btnSave := widgets.NewButton("save", pixel.V(425, 29), pixel.V(60, 30))

	textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))
	statusText := text.New(pixel.V(0, 0), textAtlas)

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
			wd.world.Paused = false
		}
		if btnStop.Draw(win) {
			wd.world.Paused = true
		}

		if btnNormal.Draw(win) {
			wd.Filter = W_FILTER_DISABLE
		}
		if btnHealth.Draw(win) {
			wd.Filter = W_FILTER_HEALTH
		}
		if btnEnergy.Draw(win) {
			wd.Filter = W_FILTER_ENERGY
		}
		if btnAge.Draw(win) {
			wd.Filter = W_FILTER_AGE
		}

		if btnSave.Draw(win) {
			wd.world.Paused = true
			saveWorld(wd.world)
			wd.world.Paused = false
		}

		win.Update()
		<-ticker.C
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
	txt.Dot = pixel.V(height/2, height/2-txt.BoundsOf("A").H()/4).Floor()
	fmt.Fprintf(txt, "FPS: % 5d | Population: % 8d | Day: % 8d | Year: % 8d | Epoch: % 8d", world.Framerate, len(world.Objects), world.Ticks, world.Year, world.Epoch)
}

func saveWorld(w *world.World) {
	filename := filepicker.SaveFile("Save world", "world.gdw")
	if filename != "" {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println("Error open file to save: ", err)
			return
		}
		defer f.Close()
		worldsaver.Save(w, f)
		fmt.Println("World saved at: ", filename)
	}
}
