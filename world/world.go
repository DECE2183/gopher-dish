package world

import (
	"gopher-dish/object"
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	WorldTicksPerYear  = 10
	WorldYearsPerEpoch = 100

	WorldSunlightMultiplier = 1.0
	WorldSunlightBeginValue = 12.0
	WorldSunlightEndValue   = 0.0
	WorldSunlightBeginPos   = 0.0
	WorldSunlightEndPos     = 0.85

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

	Paused    bool
	Framerate uint

	Ticks uint64
	Year  uint64
	Epoch uint64

	Sunlight [][]byte
	Minerals [][]byte

	Objects          map[uint64]object.Movable
	ObjectsIdCounter uint64

	Places        [][]object.Movable
	PlacesDrawMux sync.Mutex

	ticker          *time.Ticker
	state           uint
	objectsToRemove chan uint64
	chunkCount      int
	objPerChunk     int

	lastTickTime time.Time
}

func New(width, height uint32, tickInterval time.Duration) *World {
	w := &World{Width: width, Height: height}

	w.ticker = time.NewTicker(tickInterval)

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

func (w *World) ReserveID() uint64 {
	w.ObjectsIdCounter++
	return w.ObjectsIdCounter
}

func (w *World) RemoveObject(id uint64) {
	switch w.state {
	case WORLD_STATE_PREPARE:
		w.removeObject(id)
	case WORLD_STATE_HANDLE:
		w.queueObjectRemoval(id)
	}
}

func (w *World) PlaceObject(obj object.Movable, pos object.Position) bool {
	if w.state != WORLD_STATE_PREPARE && w.state != WORLD_STATE_NONE {
		return false
	}

	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return false
	}

	if !w.IsPlaceFree(pos) {
		return false
	}

	w.Objects[obj.GetID()] = obj
	w.Places[pos.X][pos.Y] = obj

	return true
}

func (w *World) MoveObject(obj object.Movable, pos object.Position) bool {
	if w.state != WORLD_STATE_PREPARE {
		return false
	}

	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return false
	}

	if !w.IsPlaceFree(pos) {
		return false
	}

	currentPos := obj.GetPosition()
	w.Places[currentPos.X][currentPos.Y] = nil
	w.Places[pos.X][pos.Y] = obj

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
	if w.state != WORLD_STATE_PREPARE && w.state != WORLD_STATE_NONE {
		return false
	}

	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return false
	}

	return (w.Places[pos.X][pos.Y] == nil)
}

func (w *World) GetCenter() object.Position {
	return object.Position{X: int32(w.Width / 2), Y: int32(w.Height / 2)}
}

func (w *World) GetFreePlace(pos object.Position) (object.Position, bool) {
	if w.state != WORLD_STATE_PREPARE {
		return object.Position{}, false
	}

	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
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
	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return 0
	}
	return w.Sunlight[pos.X][pos.Y]
}

func (w *World) GetMineralsAtPosition(pos object.Position) byte {
	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return 0
	}
	return w.Minerals[pos.X][pos.Y]
}

func (w *World) GetObjectAtPosition(pos object.Position) object.Movable {
	if w.state != WORLD_STATE_PREPARE {
		return nil
	}

	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return nil
	}

	return w.Places[pos.X][pos.Y]
}

func (w *World) SetTickPeriod(t time.Duration) {
	if t > 0 {
		w.ticker.Reset(t)
	}
}

func (w *World) Handle() {
	if w.Paused {
		time.Sleep(100 * time.Millisecond)
		w.lastTickTime = time.Now()
		return
	}

	var yearChanged, epochChanged bool

	w.Ticks++
	if w.Ticks%WorldTicksPerYear == 0 {
		yearChanged = true
		w.Year++

		if w.Year%WorldYearsPerEpoch == 0 {
			epochChanged = true
			w.Epoch++
		}
	}

	w.PlacesDrawMux.Lock()

	w.state = WORLD_STATE_PREPARE
	for _, o := range w.Objects {
		o.Prepare()
	}

	var wg sync.WaitGroup
	wg.Add(len(w.Objects))
	w.objectsToRemove = make(chan uint64, w.Width*w.Height)
	w.state = WORLD_STATE_HANDLE

	for _, o := range w.Objects {
		go func(obj object.Movable) {
			obj.Handle(yearChanged, epochChanged)
			wg.Done()
		}(o)
	}

	wg.Wait()
	close(w.objectsToRemove)
	removedObjects := 0

	w.state = WORLD_STATE_PREPARE

	for id := range w.objectsToRemove {
		w.removeObject(id)
		removedObjects++
	}

	w.PlacesDrawMux.Unlock()

	<-w.ticker.C
	w.Framerate = uint(1000000 / (time.Since(w.lastTickTime).Microseconds() + 1))
	w.lastTickTime = time.Now()
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
	if w.state != WORLD_STATE_PREPARE {
		return
	}

	obj, found := w.Objects[id]
	if !found {
		return
	}

	pos := obj.GetPosition()

	pos.X = (pos.X + int32(w.Width)) % int32(w.Width)
	if pos.Y < 0 || pos.Y >= int32(w.Height) {
		return
	}

	w.Places[pos.X][pos.Y] = nil
	delete(w.Objects, id)
}

func (w *World) queueObjectRemoval(id uint64) {
	w.objectsToRemove <- id
}
