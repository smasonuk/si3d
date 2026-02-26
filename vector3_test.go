package gosie3d

import (
	"math"
	"testing"
)

func TestNewVector3(t *testing.T) {
	v := NewVector3(1, 2, 3)
	if v.X != 1 || v.Y != 2 || v.Z != 3 || v.W != 1.0 {
		t.Errorf("NewVector3(1, 2, 3) = %v; want X=1, Y=2, Z=3, W=1.0", v)
	}
}

func TestNewVector3Full(t *testing.T) {
	v := NewVector3Full(1, 2, 3, 4)
	if v.X != 1 || v.Y != 2 || v.Z != 3 || v.W != 4.0 {
		t.Errorf("NewVector3Full(1, 2, 3, 4) = %v; want X=1, Y=2, Z=3, W=4.0", v)
	}
}

func TestNewVector3dFromArray(t *testing.T) {
	arr := []float64{1, 2, 3}
	v := NewVector3dFromArray(arr)
	if v.X != 1 || v.Y != 2 || v.Z != 3 {
		t.Errorf("NewVector3dFromArray(%v) = %v; want X=1, Y=2, Z=3", arr, v)
	}
}

func TestVector3_Add(t *testing.T) {
	v := NewVector3(1, 2, 3)
	v = v.Add(NewVector3(4, 5, 6))
	if v.X != 5 || v.Y != 7 || v.Z != 9 {
		t.Errorf("v.Add(4, 5, 6) resulted in %v; want X=5, Y=7, Z=9", v)
	}
}

func TestVector3_Normalize(t *testing.T) {
	v := NewVector3(3, 4, 0)
	v = v.Normalize()
	if math.Abs(v.X-0.6) > 1e-9 || math.Abs(v.Y-0.8) > 1e-9 || v.Z != 0 {
		t.Errorf("v.Normalize() resulted in %v; want X=0.6, Y=0.8, Z=0", v)
	}

	v2 := NewVector3(0, 0, 0)
	v2 = v2.Normalize()
	if v2.X != 0 || v2.Y != 0 || v2.Z != 0 {
		t.Errorf("v2.Normalize() resulted in %v; want X=0, Y=0, Z=0", v2)
	}
}

func TestVector3_Copy(t *testing.T) {
	v := NewVector3(1, 2, 3)
	v2 := v.Copy()
	// v == v2 comparison for values compares fields.
	if v != v2 {
		t.Errorf("v.Copy() resulted in %v; want same values as %v", v2, v)
	}
}

func TestVector3_DistanceTo(t *testing.T) {
	v1 := NewVector3(1, 1, 1)
	v2 := NewVector3(4, 5, 1)
	dist := v1.DistanceTo(v2)
	// dx = 3, dy = 4, dz = 0 -> sqrt(9+16) = 5
	if math.Abs(dist-5.0) > 1e-9 {
		t.Errorf("v1.DistanceTo(v2) = %f; want 5.0", dist)
	}
}

func TestSubtract(t *testing.T) {
	v1 := NewVector3(5, 7, 9)
	v2 := NewVector3(1, 2, 3)
	v3 := Subtract(v1, v2)
	if v3.X != 4 || v3.Y != 5 || v3.Z != 6 {
		t.Errorf("Subtract(v1, v2) = %v; want X=4, Y=5, Z=6", v3)
	}
}

func TestCross(t *testing.T) {
	v1 := NewVector3(1, 0, 0)
	v2 := NewVector3(0, 1, 0)
	v3 := Cross(v1, v2)
	// Cross product of X and Y axis is Z axis
	if v3.X != 0 || v3.Y != 0 || v3.Z != 1 {
		t.Errorf("Cross(v1, v2) = %v; want X=0, Y=0, Z=1", v3)
	}

	v4 := Cross(v2, v1)
	if v4.X != 0 || v4.Y != 0 || v4.Z != -1 {
		t.Errorf("Cross(v2, v1) = %v; want X=0, Y=0, Z=-1", v4)
	}
}

func TestDot(t *testing.T) {
	v1 := NewVector3(1, 2, 3)
	v2 := NewVector3(4, 5, 6)
	dot := Dot(v1, v2)
	// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	if dot != 32 {
		t.Errorf("Dot(v1, v2) = %f; want 32", dot)
	}
}
