package world

import (
	"gopher-dish/object"
	"math"
	"runtime"
	"sync"
)

const (
	WorldTicksPerYear  = 10
	WorldYearsPerEpoch = 100

	WorldSunlightMultiplier = 1.0
	WorldSunlightBeginValue = 25.0
	WorldSunlightEndValue   = 0.0
	WorldSunlightBeginPos   = 0.0
	WorldSunlightEndPos     = 0.6

	WorldMineralsMultiplier = 1.0
	WorldMineralsBeginValue = 1.0
	WorldMineralsEndValue   = 4.0
	WorldMineralsBeginPos   = 0.4
	WorldMineralsEndPos     = 1.0
)

const (
	WORLD_STATE_NONE = iota
	WORLD_STATE_PREPARE
	WORLD_STATE_HANDLE
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

	Objects          map[uint64]object.Movable
	ObjectsIdCounter uint64

	Places [][]object.Movable

	state           uint
	objectsToRemove chan uint64
	chunkCount      int
	objPerChunk     int
}

func New(width, height uint32) *World {
	w := &World{Width: width, Height: height}

	w.Places = make([][]object.Movable, w.Width)
	for i := 0; i < int(w.Width); i++ {
		w.Places[i] = make([]object.Movable, w.Height)
	}

	w.Objects = make(map[uint64]object.Movable)

	w.calculateSunlight()
	w.calculateMinerals()

	w.chunkCount = runtime.NumCPU()

	return w
}

func (w *World) AddObject(obj object.Movable) (uint64, object.Position) {
	pos := obj.GetPosition()

	if !w.IsPlaceFree(pos) {
		var found bool
		pos, found = w.GetFreePlace(pos)
		if !found {
			return 0, pos
		}
	}

	w.ObjectsIdCounter++
	w.Objects[w.ObjectsIdCounter] = obj
	w.Places[pos.X][pos.Y] = obj

	return w.ObjectsIdCounter, pos
}

func (w *World) RemoveObject(id uint64) {
	switch w.state {
	case WORLD_STATE_PREPARE:
		w.removeObject(id)
	case WORLD_STATE_HANDLE:
		w.queueObjectRemoval(id)
	}
}

func (w *World) MoveObject(obj object.Movable, pos object.Position) bool {
	if pos.X < 0 || pos.X >= int32(w.Width) {
		return false
	}
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return false
	}

	if !w.IsPlaceFree(pos) {
		return false
	}

	currentPos := obj.GetPosition()
	w.Places[currentPos.X][currentPos.Y] = nil
	w.Places[pos.X][pos.Y] = obj
	obj.MoveToPosition(pos)

	return true
}

func (w *World) GetObject(id uint64) object.Movable {
	obj, exists := w.Objects[id]

	if exists {
		return obj
	} else {
		return nil
	}
}

func (w *World) IsPlaceFree(pos object.Position) bool {
	return (w.Places[pos.X][pos.Y] == nil)
}

func (w *World) GetCenter() object.Position {
	return object.Position{X: int32(w.Width / 2), Y: int32(w.Height / 2)}
}

func (w *World) GetFreePlace(pos object.Position) (object.Position, bool) {
	if pos.X < 0 || pos.X >= int32(w.Width) {
		return object.Position{}, false
	}
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return object.Position{}, false
	}

	var fromX, fromY int32 = pos.X - 1, pos.Y - 1
	var toX, toY int32 = fromX + 3, fromY + 3

	if fromX < 0 {
		fromX = 0
	} else if toX >= int32(w.Width) {
		toX = int32(w.Width - 1)
	}

	if fromY < 0 {
		fromY = 0
	} else if toY >= int32(w.Height) {
		toY = int32(w.Height - 1)
	}

	for x := fromX; x < toX; x++ {
		for y := fromY; y < toY; y++ {
			if w.Places[x][y] == nil {
				return object.Position{X: x, Y: y}, true
			}
		}
	}

	return object.Position{}, false
}

func (w *World) GetSunlightAtPosition(pos object.Position) byte {
	if pos.X < 0 || pos.X >= int32(w.Width) {
		return 0
	}
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return 0
	}
	return w.Sunlight[pos.X][pos.Y]
}

func (w *World) GetMineralsAtPosition(pos object.Position) byte {
	if pos.X < 0 || pos.X >= int32(w.Width) {
		return 0
	}
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return 0
	}
	return w.Minerals[pos.X][pos.Y]
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

	w.state = WORLD_STATE_PREPARE
	for _, o := range w.Objects {
		o.Prepare()
	}

	// objIndex := 0
	// for chunk := 0; chunk < w.chunkCount; chunk++ {
	// 	for ; objIndex < w.objPerChunk; objIndex++ {
	// 		go func(o object.Movable) {

	// 		}(w.Objects[uint64(objIndex)])
	// 	}
	// }

	var wg sync.WaitGroup
	wg.Add(len(w.Objects))
	w.objectsToRemove = make(chan uint64, w.Width*w.Height)
	w.state = WORLD_STATE_HANDLE

	for _, o := range w.Objects {
		go func(obj object.Movable) {
			obj.Handle(yearChanged, epochChanged)
			wg.Done()
		}(o)

		// o.Handle(yearChanged, epochChanged)
	}

	wg.Wait()
	close(w.objectsToRemove)

	for id := range w.objectsToRemove {
		w.removeObject(id)
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

func (w *World) removeObject(id uint64) {
	obj, found := w.Objects[id]
	if !found {
		return
	}

	pos := obj.GetPosition()
	w.Places[pos.X][pos.Y] = nil

	delete(w.Objects, id)
}

func (w *World) queueObjectRemoval(id uint64) {
	w.objectsToRemove <- id
}
