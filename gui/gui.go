package gui

import (
	"time"

	"gopher-dish/object"
	"gopher-dish/world"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	updateTimerInterval time.Duration
	worldToDraw         *world.World
)

func Run(updateInterval time.Duration, world *world.World) {
	updateTimerInterval = updateInterval
	worldToDraw = world

	pixelgl.Run(initDone)
}

func initDone() {
	worldDrawer := imdraw.New(nil)
	objectsDrawer := imdraw.New(nil)
	worldBounds := pixel.R(0, 0, 1024, 480)

	drawWorld(worldDrawer, worldToDraw, worldBounds)

	cfg := pixelgl.WindowConfig{
		Title:  "gopher-dish",
		Bounds: worldBounds,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.White)
		worldDrawer.Draw(win)
		objectsDrawer.Clear()
		drawObjects(objectsDrawer, worldToDraw, worldBounds)
		objectsDrawer.Draw(win)
		win.Update()
	}
}

func drawWorld(imd *imdraw.IMDraw, world *world.World, bounds pixel.Rect) {
	sunlightColor := pixel.RGB(1, 0.83, 0.3)
	mineralsColor := pixel.RGB(0.3, 0.74, 1)

	for x := int32(0); x < int32(world.Width); x++ {
		for y := int32(0); y < int32(world.Height); y++ {
			sunlight := world.GetSunlightAtPosition(object.Position{x, y})
			minerals := world.GetMineralsAtPosition(object.Position{x, y})

			pixelColor := sunlightColor.Mul(pixel.Alpha(float64(sunlight) / 85))
			pixelColor = pixelColor.Add(mineralsColor.Mul(pixel.Alpha(float64(minerals) / 85)))

			imd.Color = pixelColor

			var posX, posY = float64(x * 6), bounds.Max.Y - float64(y*6)
			imd.Push(pixel.V(posX, posY), pixel.V(posX+6, posY-6))
			imd.Rectangle(0)
		}
	}
}

func drawObjects(imd *imdraw.IMDraw, world *world.World, bounds pixel.Rect) {
	if <-world.PlacesUpdated != true {
		return
	}

	for x := 0; x < int(world.Width); x++ {
		for y := 0; y < int(world.Height); y++ {
			if world.Places[x][y] == nil {
				continue
			}

			var posX, posY = float64(x * 6), bounds.Max.Y - float64(y*6)
			imd.EndShape = imdraw.NoEndShape

			imd.Color = colornames.Black
			imd.Push(pixel.V(posX, posY), pixel.V(posX+6, posY-6))
			imd.Rectangle(1)

			imd.Color = colornames.Green
			imd.Push(pixel.V(posX, posY-1), pixel.V(posX+5, posY-6))
			imd.Rectangle(0)
		}
	}
}
