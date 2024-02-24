package gui

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"time"

	"gopher-dish/cell"
	"gopher-dish/gui/filepicker"
	"gopher-dish/gui/fonts"
	"gopher-dish/gui/widgets"
	"gopher-dish/object"
	"gopher-dish/world"
	"gopher-dish/world/worldsaver"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

//go:embed shaders/blur.frag
var shaderBlur string

var worldTrendName = map[world.WorldEpochTrend]string{
	world.TREND_NORMAL: "normal",
	world.TREND_WARM:   "warm",
	world.TREND_COLD:   "cold",
	world.TREND_COUNT:  "unknown",
}

var (
	ticker      *time.Ticker
	worldToDraw *world.World
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

	btnPlay := widgets.NewButton("play", pixel.V(win.Bounds().W()-310, win.Bounds().H()-46), pixel.V(300, 36))
	btnStop := widgets.NewButton("stop", pixel.V(70, 29), pixel.V(60, 30))

	btnNormal := widgets.NewButton("normal", pixel.V(150, 29), pixel.V(60, 30))
	btnHealth := widgets.NewButton("health", pixel.V(215, 29), pixel.V(60, 30))
	btnEnergy := widgets.NewButton("energy", pixel.V(280, 29), pixel.V(60, 30))
	btnAge := widgets.NewButton("age", pixel.V(345, 29), pixel.V(60, 30))

	btnSave := widgets.NewButton("save", pixel.V(425, 29), pixel.V(60, 30))
	btnRestart := widgets.NewButton("restart", pixel.V(490, 29), pixel.V(60, 30))

	statusText := text.New(pixel.V(0, 0), fonts.RedhatMonoMedium12)

	viewCanvas := pixelgl.NewCanvas(pixel.R(0, 0, cfg.Bounds.Max.X, cfg.Bounds.Max.Y-24))
	viewCanvas.SetSmooth(true)

	sidePanelCanvas := pixelgl.NewCanvas(pixel.R(0, 0, 320, cfg.Bounds.Max.Y-24))
	sidePanelCanvas.SetFragmentShader(shaderBlur)
	sidePanelCanvas.SetSmooth(true)

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

		viewCanvas.SetBounds(pixel.R(0, 0, win.Bounds().W(), win.Bounds().H()-24))
		viewCanvas.Clear(colornames.White)
		wd.Draw(viewCanvas)
		viewCanvas.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().W()/2, win.Bounds().H()/2)))

		sidePanelCanvas.SetBounds(pixel.R(0, 0, 320, win.Bounds().H()-24))
		sidePanelCanvas.Clear(colornames.White)
		viewCanvas.Draw(sidePanelCanvas, pixel.IM.Moved(pixel.V(sidePanelCanvas.Bounds().W()-viewCanvas.Bounds().W()/2, viewCanvas.Bounds().H()/2)))
		sidePanelCanvas.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().W()-160, win.Bounds().H()/2+12)))

		createStatusBar(statusPanelBg, win.Bounds().Max.X, 24)
		printStatus(statusText, wd.world, win.Bounds().Max.X, 24)
		statusPanelBg.Draw(win)
		statusText.Draw(win, pixel.IM)

		btnPlay.SetPos(pixel.V(win.Bounds().W()-310, win.Bounds().H()-46))

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
		if btnRestart.Draw(win) {
			wd.world.Paused = true
			time.Sleep(time.Millisecond * 100)
			for id := range wd.world.Objects {
				wd.world.RemoveObject(id)
			}
			pos := object.Position{}
			for x := 0; x < int(wd.world.Width); x += 8 {
				for y := 0; y < int(wd.world.Height)/8; y += 1 {
					pos.X = int32(x) + 4
					pos.Y = int32(y*8) + 4*int32((x/8)%2)
					c := cell.New(wd.world, nil, pos)
					if c == nil {
						continue
					}
					c.Rotation.Degree = int32((rand.Uint32() % 8) * 45)
					for i := 0; i < 256; i++ {
						c.Genome.Mutate()
					}
				}
			}
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

func printStatus(txt *text.Text, w *world.World, width, height float64) {
	if w.Ticks%10 != 0 {
		return
	}

	txt.Clear()
	txt.Color = pixel.RGB(0.53, 0.53, 0.53)
	txt.Dot = pixel.V(height/2, height/2-txt.BoundsOf("A").H()/4).Floor()
	fmt.Fprintf(txt, "FPS: % 5d | Population: % 8d | Day: % 8d | Year: % 8d | Epoch: % 8d | Trend: % 8s",
		w.Framerate, len(w.Objects), w.Ticks, w.Year, w.Epoch, worldTrendName[w.Trend])
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
