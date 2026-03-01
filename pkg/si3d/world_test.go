package si3d

import (
	"image"
	"image/color"
	"testing"
)

// MockBatcher is a minimal PolygonBatcher implementation for testing.
type MockBatcher struct {
	drawn bool
}

func (m *MockBatcher) AddPolygon(xp, yp []float32, clr color.RGBA)                                           {}
func (m *MockBatcher) AddPolygonOutline(xp, yp []float32, strokeWidth float32, clr color.RGBA)               {}
func (m *MockBatcher) AddPolygonAndOutline(xp, yp []float32, fillClr, strokeClr color.RGBA, sw float32)      {}
func (m *MockBatcher) Draw(target *image.RGBA)                                                                { m.drawn = true }
func (m *MockBatcher) ClipPolygon(subjectPolygon []Point, screenWidth, screenHeight float32) []Point {
	return subjectPolygon
}

func TestNewWorld3d(t *testing.T) {
	w := NewWorld3d()
	if w == nil {
		t.Fatal("NewWorld3d returned nil")
	}
	if w.currentCamera != -1 {
		t.Errorf("NewWorld3d currentCamera expected -1, got %d", w.currentCamera)
	}
	if w.batcher == nil {
		t.Error("NewWorld3d batcher is nil")
	}
}

func TestWorld_AddObject(t *testing.T) {
	w := NewWorld3d()
	obj := NewModel()
	w.AddObject(&Entity{Model: obj, X: 1, Y: 2, Z: 3})

	if len(w.entities) != 1 {
		t.Errorf("World AddObject count expected 1, got %d", len(w.entities))
	}
	if w.entities[0].Model != obj {
		t.Error("World AddObject object mismatch")
	}
	if w.entities[0].X != 1 || w.entities[0].Y != 2 || w.entities[0].Z != 3 {
		t.Errorf("World AddObject position incorrect")
	}
}

func TestWorld_AddObjectDrawFirst(t *testing.T) {
	w := NewWorld3d()
	obj := NewModel()
	w.AddObjectDrawFirst(&Entity{Model: obj, X: 4, Y: 5, Z: 6})

	if len(w.entitiesDrawFirst) != 1 {
		t.Errorf("World AddObjectDrawFirst count expected 1, got %d", len(w.entitiesDrawFirst))
	}
	if w.entitiesDrawFirst[0].Model != obj {
		t.Error("World AddObjectDrawFirst object mismatch")
	}
	if w.entitiesDrawFirst[0].X != 4 || w.entitiesDrawFirst[0].Y != 5 || w.entitiesDrawFirst[0].Z != 6 {
		t.Errorf("World AddObjectDrawFirst position incorrect")
	}
}

func TestWorld_AddObjectDrawLast(t *testing.T) {
	w := NewWorld3d()
	obj := NewModel()
	w.AddObjectDrawLast(&Entity{Model: obj, X: 7, Y: 8, Z: 9})

	if len(w.entitiesDrawLast) != 1 {
		t.Errorf("World AddObjectDrawLast count expected 1, got %d", len(w.entitiesDrawLast))
	}
	if w.entitiesDrawLast[0].Model != obj {
		t.Error("World AddObjectDrawLast object mismatch")
	}
	if w.entitiesDrawLast[0].X != 7 || w.entitiesDrawLast[0].Y != 8 || w.entitiesDrawLast[0].Z != 9 {
		t.Errorf("World AddObjectDrawLast position incorrect")
	}
}

func TestWorld_AddCamera(t *testing.T) {
	w := NewWorld3d()
	cam := NewCamera(0, 0, 0, 0, 0, 0)
	w.AddCamera(cam, 10, 20, 30)

	if len(w.cameras) != 1 {
		t.Errorf("World AddCamera count expected 1, got %d", len(w.cameras))
	}
	if w.cameras[0] != cam {
		t.Error("World AddCamera camera mismatch")
	}
	if w.currentCamera != 0 {
		t.Errorf("World AddCamera currentCamera expected 0, got %d", w.currentCamera)
	}

	pos := w.cameras[0].GetPosition()
	if pos.X != 10 || pos.Y != 20 || pos.Z != 30 {
		t.Errorf("World AddCamera position incorrect: %v", pos)
	}

	// Add another camera
	cam2 := NewCamera(0, 0, 0, 0, 0, 0)
	w.AddCamera(cam2, 40, 50, 60)
	if w.currentCamera != 1 {
		t.Errorf("World AddCamera currentCamera expected 1, got %d", w.currentCamera)
	}
}

func TestWorld_SetPolygonBatcher(t *testing.T) {
	w := NewWorld3d()
	mock := &MockBatcher{}
	w.SetPolygonBatcher(mock)

	// Build a minimal scene: one camera and one compiled model.
	cam := NewCamera(0, 0, 0, 0, 0, 0)
	w.AddCamera(cam, 0, 0, -50)

	model := NewCube()
	w.AddObject(&Entity{Model: model, X: 0, Y: 0, Z: 0})

	target := image.NewRGBA(image.Rect(0, 0, 320, 240))
	w.PaintObjects(target, 320, 240)

	if !mock.drawn {
		t.Error("SetPolygonBatcher: expected MockBatcher.Draw to be called during PaintObjects")
	}
}
