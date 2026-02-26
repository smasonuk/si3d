package si3d

import (
	"image/color"
	"testing"
)

func TestNewFaceEmpty(t *testing.T) {
	col := color.RGBA{255, 0, 0, 255}
	norm := NewVector3(0, 1, 0)
	f := NewFaceEmpty(col, norm)

	if f == nil {
		t.Fatal("NewFaceEmpty returned nil")
	}
	if len(f.Points) != 0 {
		t.Errorf("NewFaceEmpty should have 0 points, got %d", len(f.Points))
	}
	if f.Col != col {
		t.Errorf("NewFaceEmpty color mismatch")
	}
	if f.normal.X != norm.X || f.normal.Y != norm.Y || f.normal.Z != norm.Z {
		t.Errorf("NewFaceEmpty normal mismatch")
	}
}

func TestNewFace(t *testing.T) {
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(1, 0, 0)
	p3 := NewVector3(0, 1, 0)
	pnts := []Vector3{p1, p2, p3}
	col := color.RGBA{0, 255, 0, 255}
	norm := NewVector3(0, 0, 1)

	f := NewFace(pnts, col, norm)

	if f == nil {
		t.Fatal("NewFace returned nil")
	}
	if f.Cnum != 3 {
		t.Errorf("NewFace Cnum should be 3, got %d", f.Cnum)
	}
	if len(f.Points) != 3 {
		t.Errorf("NewFace Points length should be 3, got %d", len(f.Points))
	}
}

func TestFace_Copy(t *testing.T) {
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(1, 0, 0)
	p3 := NewVector3(0, 1, 0)
	pnts := []Vector3{p1, p2, p3}
	col := color.RGBA{0, 0, 255, 255}
	norm := NewVector3(0, 0, 1)

	f := NewFace(pnts, col, norm)
	// Make sure normal is initialized
	f.GetNormal()

	fCopy := f.Copy()

	if fCopy == nil {
		t.Fatal("Copy returned nil")
	}
	if len(fCopy.Points) != len(f.Points) {
		t.Errorf("Copy Points length mismatch")
	}
	if fCopy.Col != f.Col {
		t.Errorf("Copy color mismatch")
	}
	if fCopy.GetNormal().X != f.GetNormal().X {
		t.Errorf("Copy normal mismatch")
	}

	// Verify deep copy of points
	f.Points[0].X = 100
	if fCopy.Points[0].X == 100 {
		t.Errorf("Copy should be deep copy for Points")
	}
}

func TestFace_AddPoint_And_Finished(t *testing.T) {
	col := color.RGBA{255, 255, 0, 255}
	norm := NewVector3(0, 0, 1)
	f := NewFaceEmpty(col, norm)

	f.AddPoint(0, 0, 0)
	f.AddPoint(1, 0, 0)
	f.AddPoint(0, 1, 0)

	if len(f.vecPnts) != 3 {
		t.Errorf("AddPoint failed, vecPnts length should be 3, got %d", len(f.vecPnts))
	}

	f.Finished(FACE_NORMAL)

	if len(f.Points) != 3 {
		t.Errorf("Finished failed, Points length should be 3, got %d", len(f.Points))
	}
	if f.Cnum != 3 {
		t.Errorf("Finished failed, Cnum should be 3, got %d", f.Cnum)
	}
	if f.vecPnts != nil {
		t.Errorf("Finished failed, vecPnts should be nil")
	}
}

func TestFace_GetNormal(t *testing.T) {
	// Triangle in XY plane, normal should be Z
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(1, 0, 0)
	p3 := NewVector3(0, 1, 0)
	pnts := []Vector3{p1, p2, p3}
	col := color.RGBA{255, 255, 255, 255}

	// NewFace with nil normal to force calculation
	f := NewFace(pnts, col, Vector3{})

	n := f.GetNormal()
	if n.X != 0 || n.Y != 0 || n.Z != 1 {
		t.Errorf("GetNormal calculated incorrect normal: %v, expected (0,0,1)", n)
	}

	// Test with explicit normal set
	explicitNormal := NewVector3(1, 1, 1)
	f.SetNormal(explicitNormal)
	n2 := f.GetNormal()
	if n2.X != 1 || n2.Y != 1 || n2.Z != 1 {
		t.Errorf("SetNormal/GetNormal failed: %v", n2)
	}
}

func TestFace_GetPlane(t *testing.T) {
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(1, 0, 0)
	p3 := NewVector3(0, 1, 0)
	pnts := []Vector3{p1, p2, p3}
	col := color.RGBA{255, 255, 255, 255}
	f := NewFace(pnts, col, Vector3{})

	plane := f.GetPlane()
	if plane == nil {
		t.Fatal("GetPlane returned nil")
	}
	// Normal is (0,0,1), point is (0,0,0) -> 0*x + 0*y + 1*z + D = 0 => D=0
	// Plane equation: z = 0
	if plane.C != 1 {
		t.Errorf("Plane normal Z component should be 1, got %f", plane.C)
	}
	if plane.D != 0 {
		t.Errorf("Plane D component should be 0, got %f", plane.D)
	}
}

func TestFace_GetMidPoint(t *testing.T) {
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(2, 0, 0)
	pnts := []Vector3{p1, p2} // Line segment
	col := color.RGBA{255, 255, 255, 255}
	f := NewFace(pnts, col, Vector3{})

	mid := f.GetMidPoint()
	if mid.X != 1 || mid.Y != 0 || mid.Z != 0 {
		t.Errorf("GetMidPoint incorrect: %v", mid)
	}
}

func TestFace_GetDistanceToPoint(t *testing.T) {
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(2, 0, 0)
	pnts := []Vector3{p1, p2} // Midpoint is (1,0,0)
	col := color.RGBA{255, 255, 255, 255}
	f := NewFace(pnts, col, Vector3{})

	pt := NewVector3(1, 1, 0) // Distance to midpoint (1,0,0) is 1
	dist := f.GetDistanceToPoint(pt)

	if dist != 1 {
		t.Errorf("GetDistanceToPoint incorrect: %f, expected 1", dist)
	}
}
