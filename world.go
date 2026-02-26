package gosie3d

import (
	"image"
	"sort"
)

const UP_DIR = -1.0

type RenderContext struct {
	Buffer3D     []Vector3
	BufferPoints []Point
	BufferFloatX []float32
	BufferFloatY []float32
}

func NewRenderContext() *RenderContext {
	return &RenderContext{
		Buffer3D:     make([]Vector3, 0, 5000),
		BufferPoints: make([]Point, 0, 100),
		BufferFloatX: make([]float32, 0, 100),
		BufferFloatY: make([]float32, 0, 100),
	}
}

type Entity struct {
	Model   *Model
	X, Y, Z float64
}

type World struct {
	entities      []*Entity
	cameras       []*Camera
	currentCamera int
	// draw first
	entitiesDrawFirst []*Entity
	// draw last
	entitiesDrawLast []*Entity
	batcher          *PolygonBatcher
	ctx              *RenderContext
}

func NewWorld3d() *World {
	return &World{
		currentCamera: -1,
		batcher:       NewPolygonBatcher(5000),
		ctx:           NewRenderContext(),
	}
}

func (w *World) AddObject(e *Entity) {
	w.entities = append(w.entities, e)
}

// AddObjectDrawFirst
func (w *World) AddObjectDrawFirst(e *Entity) {
	w.entitiesDrawFirst = append(w.entitiesDrawFirst, e)
}

// AddObjectDrawLast
func (w *World) AddObjectDrawLast(e *Entity) {
	w.entitiesDrawLast = append(w.entitiesDrawLast, e)
}

func (w *World) AddCamera(c *Camera, x, y, z float64) {
	c.SetCameraPosition(x, y, z)
	w.cameras = append(w.cameras, c)
	w.currentCamera = len(w.cameras) - 1
}

func paint(batcher *PolygonBatcher, xsize, ysize int, e *Entity, cam *Camera, ctx *RenderContext) {
	objToWorld := TransMatrix(e.X, e.Y, e.Z)

	objToCam := cam.camMatrixRev.MultiplyBy(objToWorld)
	e.Model.ApplyMatrixTemp(objToCam)
	e.Model.PaintObject(batcher, xsize/2, ysize/2, true, float32(xsize), float32(ysize), cam.GetNearPlane(), ctx)
}

func distBetweenEntityAndCamera(e *Entity, cam *Camera) float64 {
	objX, objY, objZ := e.X, e.Y, e.Z
	pos := cam.GetPosition()
	camX, camY, camZ := pos.X, pos.Y, pos.Z

	return (objX-camX)*(objX-camX) + (objY-camY)*(objY-camY) + (objZ-camZ)*(objZ-camZ)
}

func draw(batcher *PolygonBatcher, xsize, ysize int, entities []*Entity, cam *Camera, ctx *RenderContext) {
	// draw background objects
	for _, e := range entities {
		objToWorld := TransMatrix(
			e.X,
			e.Y,
			e.Z,
		)
		objToCam := cam.camMatrixRev.MultiplyBy(objToWorld)
		e.Model.ApplyMatrixTemp(objToCam)
		e.Model.PaintObject(batcher, xsize/2, ysize/2, true, float32(xsize), float32(ysize), cam.GetNearPlane(), ctx)
	}

}

func sortObjects(entities []*Entity, cam *Camera) {
	type sortableEntity struct {
		e      *Entity
		distSq float64
	}

	sortedEntities := make([]sortableEntity, len(entities))
	for i, e := range entities {
		sortedEntities[i] = sortableEntity{
			e:      e,
			distSq: distBetweenEntityAndCamera(e, cam),
		}
	}

	sort.Slice(sortedEntities, func(i, j int) bool {
		return sortedEntities[i].distSq > sortedEntities[j].distSq
	})

	for i, se := range sortedEntities {
		entities[i] = se.e
	}
}

func (w *World) PaintObjects(target *image.RGBA, xsize, ysize int) {

	if w.currentCamera == -1 || len(w.cameras) == 0 {
		return
	}
	cam := w.cameras[w.currentCamera]

	entitiesToDraw := make([]*Entity, len(w.entities))
	copy(entitiesToDraw, w.entities)

	type sortableEntity struct {
		e      *Entity
		distSq float64
	}

	sortedEntities := make([]sortableEntity, len(entitiesToDraw))
	for i, e := range entitiesToDraw {
		sortedEntities[i] = sortableEntity{
			e:      e,
			distSq: distBetweenEntityAndCamera(e, cam),
		}
	}

	sort.Slice(sortedEntities, func(i, j int) bool {
		return sortedEntities[i].distSq > sortedEntities[j].distSq
	})

	for i, se := range sortedEntities {
		entitiesToDraw[i] = se.e
	}

	// // Draw objects that should be drawn first (that dont have a direction vector)
	for _, e := range w.entitiesDrawFirst {
		if e.Model.hasObjectDirection {
			continue
		}

		paint(w.batcher, xsize, ysize, e, cam, w.ctx)
	}

	// get objects which are poing at and from the camera. objects pointing towards the camera are drawn first
	backgroundObjects := make([]*Entity, 0)
	foregroundObjects := make([]*Entity, 0)
	for _, e := range w.entitiesDrawFirst {
		if !e.Model.hasObjectDirection {
			continue
		}

		objToWorld := TransMatrix(e.X, e.Y, e.Z)
		objToCam := cam.camMatrixRev.MultiplyBy(objToWorld)

		// Transform the direction vector (the direction the cusion is pointing) to camera space
		direction := objToCam.RotateVector3(e.Model.objectDirection)
		direction = direction.Normalize()

		e.Model.ApplyMatrixTemp(objToCam)
		pointX := e.Model.transFaceMesh.Points[0].X
		pointY := e.Model.transFaceMesh.Points[0].Y
		pointZ := e.Model.transFaceMesh.Points[0].Z

		plane := NewPlaneFromPoint(NewVector3(pointX, pointY, pointZ), direction)
		where := plane.PointOnPlane(0, 0, 0)

		if where > 0 {
			backgroundObjects = append(backgroundObjects, e)
		} else {
			foregroundObjects = append(foregroundObjects, e)
		}
	}

	// sort background objects by distance to camera and draw them
	sortObjects(backgroundObjects, cam)
	draw(w.batcher, xsize, ysize, backgroundObjects, cam, w.ctx)

	for _, e := range entitiesToDraw {
		// object space to world space trans
		objToWorld := TransMatrix(
			e.X,
			e.Y,
			e.Z,
		)

		// objToCam := objToWorld.MultiplyBy(cam.camMatrixRev)
		objToCam := cam.camMatrixRev.MultiplyBy(objToWorld)

		// then, world space to camera space
		e.Model.ApplyMatrixTemp(objToCam)
		e.Model.PaintObject(w.batcher, xsize/2, ysize/2, true, float32(xsize), float32(ysize), cam.GetNearPlane(), w.ctx)
	}

	// draw foreground objects
	sortObjects(foregroundObjects, cam)
	draw(w.batcher, xsize, ysize, foregroundObjects, cam, w.ctx)

	// Draw objects that should be drawn last
	for _, e := range w.entitiesDrawLast {
		paint(w.batcher, xsize, ysize, e, cam, w.ctx)
	}

	w.batcher.Draw(target)
}
