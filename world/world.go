package world

import (
	"gopher-dish/object"
	"math"
)

const (
	WorldTicksPerYear  = 10
	WorldYearsPerEpoch = 100

	WorldSunlightMultiplier = 1.0
	WorldSunlightBeginValue = 25.0
	WorldSunlightEndValue   = 0.0
	WorldSunlightBeginPos   = 0.0
	WorldSunlightEndPos     = 0.5

	WorldMineralsMultiplier = 1.0
	WorldMineralsBeginValue = 1.0
	WorldMineralsEndValue   = 4.0
	WorldMineralsBeginPos   = 0.5
	WorldMineralsEndPos     = 1.0
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func remap[T number](value, inMin, inMax, outMin, outMax T) T {
	return outMin + (value-inMin)*(outMax-outMin)/(inMax-inMin)
}

type World struct {
	Width, Height uint32

	Ticks uint64
	Year  uint64
	Epoch uint64

	Sunlight [][]byte
	Minerals [][]byte

	Objects          map[uint64]object.Object
	ObjectsIdCounter uint64
}

func New(width, height uint32) *World {
	w := &World{Width: width, Height: height}

	w.Objects = make(map[uint64]object.Object)

	w.calculateSunlight()
	w.calculateMinerals()

	return w
}

func (w *World) AddObject(obj object.Object) uint64 {
	w.ObjectsIdCounter++
	w.Objects[w.ObjectsIdCounter] = obj
	return w.ObjectsIdCounter
}

func (w *World) GetObject(id uint64) object.Object {
	obj, exists := w.Objects[id]

	if exists {
		return obj
	} else {
		return nil
	}
}

func (w *World) GetSunlightAtPosition(pos object.Position) byte {
	if pos.X < 0 || pos.X > int32(w.Width) {
		return 0
	}
	if pos.Y < 0 || pos.Y > int32(w.Height) {
		return 0
	}
	return w.Sunlight[pos.X][pos.Y]
}

func (w *World) GetMineralsAtPosition(pos object.Position) byte {
	if pos.X < 0 || pos.X > int32(w.Width) {
		return 0
	}
	if pos.Y < 0 || pos.Y > int32(w.Height) {
		return 0
	}
	return w.Minerals[pos.X][pos.Y]
}

func (w *World) RemoveObject(id uint64) {
	delete(w.Objects, id)
}

func (w *World) Handle() {
	var yearChanged, epochChanged bool

	w.Ticks++
	if w.Ticks%WorldTicksPerYear == 0 {
		yearChanged = true
		w.Year++
	}
	if w.Year%WorldYearsPerEpoch == 0 {
		epochChanged = true
		w.Epoch++
	}

	for _, o := range w.Objects {
		o.Handle(yearChanged, epochChanged)
	}
}

func (w *World) calculateSunlight() {
	sunlightBegin := uint32(math.Round(float64(w.Height) * WorldSunlightBeginPos))
	sunlightEnd := uint32(math.Round(float64(w.Height) * WorldSunlightEndPos))

	w.Sunlight = make([][]byte, w.Width)
	for x := 0; x < int(w.Width); x++ {
		w.Sunlight[x] = make([]byte, w.Height)
		for y := sunlightBegin; y < sunlightEnd; y++ {
			sunlightValue := remap(float64(y), float64(sunlightBegin), float64(sunlightEnd), WorldSunlightBeginValue, WorldSunlightEndValue)
			w.Sunlight[x][y] = byte(math.Round(sunlightValue * WorldSunlightMultiplier))
		}
	}
}

func (w *World) calculateMinerals() {
	mineralsBegin := uint32(math.Round(float64(w.Height) * WorldMineralsBeginPos))
	mineralsEnd := uint32(math.Round(float64(w.Height) * WorldMineralsEndPos))

	w.Minerals = make([][]byte, w.Width)
	for x := 0; x < int(w.Width); x++ {
		w.Minerals[x] = make([]byte, w.Height)
		for y := mineralsBegin; y < mineralsEnd; y++ {
			mineralsCount := remap(float64(y), float64(mineralsBegin), float64(mineralsEnd), WorldMineralsBeginValue, WorldMineralsEndValue)
			w.Minerals[x][y] = byte(math.Round(mineralsCount * WorldMineralsMultiplier))
		}
	}
}
