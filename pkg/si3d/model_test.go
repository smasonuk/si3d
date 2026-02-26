package si3d

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	obj := NewModel()
	if obj == nil {
		t.Fatal("NewModel returned nil")
	}
	if obj.transFaceMesh == nil {
		t.Error("transFaceMesh is nil")
	}
	if obj.transNormalMesh == nil {
		t.Error("transNormalMesh is nil")
	}
	if obj.faces == nil {
		t.Error("faces is nil")
	}
}

func TestModel_SetGetPosition(t *testing.T) {
	obj := NewModel()
	pos := obj.GetPosition()
	if pos.X != 0 || pos.Y != 0 || pos.Z != 0 {
		t.Error("Initial position should be (0,0,0)")
	}

	obj.SetPosition(10, 20, 30)
	pos = obj.GetPosition()
	if pos.X != 10 || pos.Y != 20 || pos.Z != 30 {
		t.Errorf("SetPosition failed, got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}
}

func TestModel_Clone(t *testing.T) {
	obj := NewModel()
	// Modify something to verify deep copy of some parts or shared others
	// According to Clone implementation:
	// faceMesh, normalMesh, faces, root, faceIndices, normalIndices are shared
	// transFaceMesh, transNormalMesh are copied

	obj.transFaceMesh.AddPoint(NewVector3(1, 2, 3))

	clone := obj.Clone()

	if len(clone.transFaceMesh.Points) != len(obj.transFaceMesh.Points) {
		t.Error("Clone did not copy transFaceMesh points")
	}

	// Check that it is a deep copy for transFaceMesh
	clone.transFaceMesh.Points[0].X = 999
	if obj.transFaceMesh.Points[0].X == 999 {
		t.Error("Clone transFaceMesh is not a deep copy")
	}
}

func TestModel_TranslateAllPoints(t *testing.T) {
	obj := NewModel()
	// Need to initialize faceMesh and points because TranslateAllPoints operates on faceMesh
	obj.faceMesh = NewFaceMesh()
	obj.faceMesh.AddPoint(NewVector3(0, 0, 0))
	obj.faceMesh.AddPoint(NewVector3(1, 1, 1))

	obj.TranslateAllPoints(10, 20, 30)

	p0 := obj.faceMesh.Points[0]
	if p0.X != 10 || p0.Y != 20 || p0.Z != 30 {
		t.Errorf("TranslateAllPoints p0 failed, got (%f, %f, %f)", p0.X, p0.Y, p0.Z)
	}

	p1 := obj.faceMesh.Points[1]
	if p1.X != 11 || p1.Y != 21 || p1.Z != 31 {
		t.Errorf("TranslateAllPoints p1 failed, got (%f, %f, %f)", p1.X, p1.Y, p1.Z)
	}
}

func TestModel_ScaleAllPoints(t *testing.T) {
	obj := NewModel()
	obj.faceMesh = NewFaceMesh()
	obj.faceMesh.AddPoint(NewVector3(1, 2, 3))

	obj.ScaleAllPoints(2.0)

	p := obj.faceMesh.Points[0]
	if p.X != 2 || p.Y != 4 || p.Z != 6 {
		t.Errorf("ScaleAllPoints failed, got (%f, %f, %f)", p.X, p.Y, p.Z)
	}
}

func TestModel_CalcSize(t *testing.T) {
	obj := NewModel()
	obj.faceMesh = NewFaceMesh()
	obj.faceMesh.AddPoint(NewVector3(0, 0, 0))
	obj.faceMesh.AddPoint(NewVector3(10, 20, 30))

	obj.CalcSize()

	if obj.XLength() != 10 {
		t.Errorf("XLength expected 10, got %f", obj.XLength())
	}
	if obj.YLength() != 20 {
		t.Errorf("YLength expected 20, got %f", obj.YLength())
	}
	if obj.ZLength() != 30 {
		t.Errorf("ZLength expected 30, got %f", obj.ZLength())
	}

	xl, yl, zl := obj.GetExtents()
	if xl != 10 || yl != 20 || zl != 30 {
		t.Errorf("GetExtents failed, got (%f, %f, %f)", xl, yl, zl)
	}
}

func TestModel_Center(t *testing.T) {
	obj := NewModel()
	obj.faceMesh = NewFaceMesh()
	// Box from (0,0,0) to (10,10,10)
	// Center is (5,5,5)
	obj.faceMesh.AddPoint(NewVector3(0, 0, 0))
	obj.faceMesh.AddPoint(NewVector3(10, 10, 10))

	obj.Center()

	// Points should be moved by (-5, -5, -5)
	// 0 -> -5
	// 10 -> 5

	p0 := obj.faceMesh.Points[0]
	if p0.X != -5 || p0.Y != -5 || p0.Z != -5 {
		t.Errorf("Center p0 failed, got (%f, %f, %f)", p0.X, p0.Y, p0.Z)
	}

	p1 := obj.faceMesh.Points[1]
	if p1.X != 5 || p1.Y != 5 || p1.Z != 5 {
		t.Errorf("Center p1 failed, got (%f, %f, %f)", p1.X, p1.Y, p1.Z)
	}
}
