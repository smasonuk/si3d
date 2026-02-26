package gosie3d

import (
	"math"
	"testing"
)

func TestTransform_GetMatrix(t *testing.T) {
	// Test translation (1, 2, 3) and scale (2, 2, 2)
	tr := NewTransform()
	tr.Position = NewVector3(1, 2, 3)
	tr.Scale = NewVector3(2, 2, 2)

	m := tr.GetMatrix()

	// With Scale * Rotation * Translation (applied to row vector v)
	// v * M = v * Scale * Rotation * Translation
	// v = (1, 1, 1).
	// Scale(2) -> (2, 2, 2).
	// Rotation(0) -> (2, 2, 2).
	// Translation(1, 2, 3) -> (3, 4, 5).

	// Let's verify manually.
	// m[0][0] should be 2. (Scale X)
	// m[1][1] should be 2. (Scale Y)
	// m[2][2] should be 2. (Scale Z)
	// m[3][3] should be 1.
	// m[3][0] should be 1. (Translation X)
	// m[3][1] should be 2. (Translation Y)
	// m[3][2] should be 3. (Translation Z)

	if m.ThisMatrix[0][0] != 2.0 {
		t.Errorf("Scale X incorrect, got %f", m.ThisMatrix[0][0])
	}
	if m.ThisMatrix[1][1] != 2.0 {
		t.Errorf("Scale Y incorrect, got %f", m.ThisMatrix[1][1])
	}
	if m.ThisMatrix[2][2] != 2.0 {
		t.Errorf("Scale Z incorrect, got %f", m.ThisMatrix[2][2])
	}

	if m.ThisMatrix[3][0] != 1.0 {
		t.Errorf("Translation X incorrect, got %f", m.ThisMatrix[3][0])
	}
	if m.ThisMatrix[3][1] != 2.0 {
		t.Errorf("Translation Y incorrect, got %f", m.ThisMatrix[3][1])
	}
	if m.ThisMatrix[3][2] != 3.0 {
		t.Errorf("Translation Z incorrect, got %f", m.ThisMatrix[3][2])
	}

	// Verify transform of a point
	v := NewVector3(1, 1, 1)
	// TransformObj uses v * M logic internally.
	src := []Vector3{v}
	dest := []Vector3{{}}
	m.TransformObj(src, dest)

	expected := NewVector3(3, 4, 5) // (1*2+1, 1*2+2, 1*2+3) -> (3, 4, 5)
	if dest[0].X != expected.X || dest[0].Y != expected.Y || dest[0].Z != expected.Z {
		t.Errorf("Transformed point incorrect, got %v, expected %v", dest[0], expected)
	}

	// Test Rotation Y (90 deg)
	tr2 := NewTransform()
	tr2.Rotate(NewVector3(0, 1, 0), math.Pi/2)
	// Rotate (1, 0, 0)
	// Based on Matrix.go ROTY implementation:
	// x' = x*Cos + z*Sin
	// z' = x*-Sin + z*Cos
	// If x=1, z=0: x'=0, z'=-1.

	m2 := tr2.GetMatrix()
	src2 := []Vector3{{X: 1, Y: 0, Z: 0, W: 1}}
	dest2 := []Vector3{{}}
	m2.TransformObj(src2, dest2)

	if math.Abs(dest2[0].X) > 1e-6 {
		t.Errorf("Rotated X incorrect, got %f", dest2[0].X)
	}
	if math.Abs(dest2[0].Z-(-1.0)) > 1e-6 {
		t.Errorf("Rotated Z incorrect, got %f", dest2[0].Z)
	}
}
