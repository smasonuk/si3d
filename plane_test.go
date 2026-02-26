package gosie3d

import (
	"image/color"
	"testing"
)

func TestNewPlaneFromPoint(t *testing.T) {
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 0, 1) // Plane z = 0
	p := NewPlaneFromPoint(pt, norm)

	if p == nil {
		t.Fatal("NewPlaneFromPoint returned nil")
	}
	if p.A != 0 || p.B != 0 || p.C != 1 {
		t.Errorf("Plane normal mismatch: (%f,%f,%f)", p.A, p.B, p.C)
	}
	if p.D != 0 {
		t.Errorf("Plane D mismatch: %f", p.D)
	}

	pt2 := NewVector3(0, 0, 5)
	p2 := NewPlaneFromPoint(pt2, norm) // Plane z = 5 => z - 5 = 0 => D = -5
	if p2.D != -5 {
		t.Errorf("Plane D mismatch for point (0,0,5): %f", p2.D)
	}
}

func TestPlane_PointOnPlane(t *testing.T) {
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 0, 1)
	p := NewPlaneFromPoint(pt, norm) // z = 0

	// Point on plane
	res := p.PointOnPlane(0, 0, 0)
	if res != 0 {
		t.Errorf("PointOnPlane failed for point on plane: %f", res)
	}

	// Point above plane (positive z)
	res = p.PointOnPlane(0, 0, 1)
	if res <= 0 {
		t.Errorf("PointOnPlane failed for point above plane: %f", res)
	}

	// Point below plane (negative z)
	res = p.PointOnPlane(0, 0, -1)
	if res >= 0 {
		t.Errorf("PointOnPlane failed for point below plane: %f", res)
	}
}

func TestPlane_lIntersect(t *testing.T) {
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 0, 1)
	p := NewPlaneFromPoint(pt, norm) // z = 0

	// Line intersects plane (from z=1 to z=-1)
	p1 := NewVector3(0, 0, 1)
	p2 := NewVector3(0, 0, -1)
	if !p.lIntersect(p1, p2) {
		t.Errorf("lIntersect failed for intersecting line")
	}

	// Line does not intersect plane (from z=1 to z=2)
	p3 := NewVector3(0, 0, 1)
	p4 := NewVector3(0, 0, 2)
	if p.lIntersect(p3, p4) {
		t.Errorf("lIntersect failed for non-intersecting line")
	}

	// Line touches plane (from z=0 to z=1)
	p5 := NewVector3(0, 0, 0)
	p6 := NewVector3(0, 0, 1)
	if p.lIntersect(p5, p6) {
		t.Errorf("lIntersect failed for touching line")
	}
}

func TestPlane_LineIntersect(t *testing.T) {
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 0, 1)
	p := NewPlaneFromPoint(pt, norm) // z = 0

	// Line intersects plane at (0,0,0)
	p1 := NewVector3(0, 0, 1)
	p2 := NewVector3(0, 0, -1)
	intersect, ok := p.LineIntersect(p1, p2)
	if !ok {
		t.Fatal("LineIntersect returned false for intersecting line")
	}
	if intersect.X != 0 || intersect.Y != 0 || intersect.Z != 0 {
		t.Errorf("LineIntersect incorrect: %v", intersect)
	}

	// Line does not intersect
	p3 := NewVector3(0, 0, 1)
	p4 := NewVector3(0, 0, 2)
	_, ok2 := p.LineIntersect(p3, p4)
	if ok2 {
		t.Errorf("LineIntersect returned true for non-intersecting line")
	}
}

func TestPlane_FaceIntersect(t *testing.T) {
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 0, 1)
	p := NewPlaneFromPoint(pt, norm) // z = 0

	// Face intersects plane
	p1 := NewVector3(0, 0, 1)
	p2 := NewVector3(1, 0, 1)
	p3 := NewVector3(0, 1, -1)
	fPoints := []Vector3{p1, p2, p3}
	f := NewFace(fPoints, color.RGBA{}, Vector3{})

	if !p.FaceIntersect(f) {
		t.Errorf("FaceIntersect failed for intersecting face")
	}

	// Face does not intersect (all points above)
	p4 := NewVector3(0, 0, 1)
	p5 := NewVector3(1, 0, 1)
	p6 := NewVector3(0, 1, 1)
	fPoints2 := []Vector3{p4, p5, p6}
	f2 := NewFace(fPoints2, color.RGBA{}, Vector3{})

	if p.FaceIntersect(f2) {
		t.Errorf("FaceIntersect failed for non-intersecting face")
	}
}

func TestPlane_SplitFace(t *testing.T) {
	// Plane z=0
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 0, 1)
	p := NewPlaneFromPoint(pt, norm)

	// Triangle intersecting plane
	p1 := NewVector3(0, 1, 1)
	p2 := NewVector3(1, 0, 1)
	p3 := NewVector3(0, 0, -1)
	fPoints := []Vector3{p1, p2, p3}
	f := NewFace(fPoints, color.RGBA{}, Vector3{})

	splitFaces := p.SplitFace(f)

	if len(splitFaces) != 2 {
		t.Errorf("SplitFace should return 2 faces, got %d", len(splitFaces))
	}

	if splitFaces[0] == nil || splitFaces[1] == nil {
		t.Errorf("SplitFace returned nil faces")
	}

	// Both faces should have points.
	if len(splitFaces[0].Points) < 3 {
		t.Errorf("SplitFace[0] has too few points: %d", len(splitFaces[0].Points))
	}
	if len(splitFaces[1].Points) < 3 {
		t.Errorf("SplitFace[1] has too few points: %d", len(splitFaces[1].Points))
	}

	// Test non-intersecting face
	p4 := NewVector3(0, 1, 2)
	p5 := NewVector3(1, 0, 2)
	p6 := NewVector3(0, 0, 2)
	fPoints2 := []Vector3{p4, p5, p6}
	f2 := NewFace(fPoints2, color.RGBA{}, Vector3{})

	splitFaces2 := p.SplitFace(f2)
	if splitFaces2[0] != f2 {
		t.Errorf("SplitFace should return original face if no intersection")
	}
	if splitFaces2[1] != nil {
		t.Errorf("SplitFace[1] should be nil if no intersection")
	}
}

func TestPointOnPlane_YNormal(t *testing.T) {
	// Create a flat plane at Y=0 pointing UP.
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 1, 0)
	p := NewPlaneFromPoint(pt, norm)

	// Test PointOnPlane
	// (0, 5, 0) should return a positive value
	res := p.PointOnPlane(0, 5, 0)
	if res <= 0 {
		t.Errorf("Expected positive value for point (0,5,0), got %f", res)
	}

	// (0, -5, 0) should return a negative value
	res = p.PointOnPlane(0, -5, 0)
	if res >= 0 {
		t.Errorf("Expected negative value for point (0,-5,0), got %f", res)
	}

	// (0, 0, 0) should return 0
	res = p.PointOnPlane(0, 0, 0)
	if res != 0 {
		t.Errorf("Expected 0 for point (0,0,0), got %f", res)
	}
}

func TestLineIntersect_Vertical(t *testing.T) {
	// Plane Y=0, Normal (0,1,0)
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(0, 1, 0)
	p := NewPlaneFromPoint(pt, norm)

	// Line going straight through the plane
	p1 := NewVector3(0, 10, 0)
	p2 := NewVector3(0, -10, 0)

	intersect, ok := p.LineIntersect(p1, p2)
	if !ok {
		t.Fatal("Expected intersection, got false")
	}

	if intersect.X != 0 || intersect.Y != 0 || intersect.Z != 0 {
		t.Errorf("Expected intersection at (0,0,0), got %v", intersect)
	}
}

func TestSplitFace_Square(t *testing.T) {
	// Plane at X=0, Normal (1,0,0)
	pt := NewVector3(0, 0, 0)
	norm := NewVector3(1, 0, 0)
	p := NewPlaneFromPoint(pt, norm)

	// Square face intersecting the plane directly down the middle
	p1 := NewVector3(-10, -10, 0)
	p2 := NewVector3(10, -10, 0)
	p3 := NewVector3(10, 10, 0)
	p4 := NewVector3(-10, 10, 0)
	fPoints := []Vector3{p1, p2, p3, p4}
	f := NewFace(fPoints, color.RGBA{}, Vector3{})

	splitFaces := p.SplitFace(f)

	if len(splitFaces) != 2 {
		t.Fatalf("Expected 2 faces, got %d", len(splitFaces))
	}

	if splitFaces[0] == nil || splitFaces[1] == nil {
		t.Fatal("One of the split faces is nil")
	}

	// Verify that the square was split into two rectangular faces
	if len(splitFaces[0].Points) != 4 {
		t.Errorf("Expected split face 0 to have 4 points, got %d", len(splitFaces[0].Points))
	}
	if len(splitFaces[1].Points) != 4 {
		t.Errorf("Expected split face 1 to have 4 points, got %d", len(splitFaces[1].Points))
	}
}
