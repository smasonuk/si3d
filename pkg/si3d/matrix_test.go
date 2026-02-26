package si3d

import (
	"math"
	"testing"
)

func TestNewMatrix(t *testing.T) {
	m := NewMatrix()
	// Check dimensions implicitly by accessing elements
	if m.ThisMatrix[0][0] != 0 {
		t.Errorf("NewMatrix() [0][0] = %f; want 0", m.ThisMatrix[0][0])
	}
}

func TestNewMatrixFromData(t *testing.T) {
	data := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	m := NewMatrixFromData(data)

	if m.ThisMatrix[0][0] != 1 || m.ThisMatrix[1][3] != 8 {
		t.Error("NewMatrixFromData() content mismatch")
	}
}

func TestNewRotationMatrix(t *testing.T) {
	// Test ROTZ 90 degrees (PI/2)
	// cos(90) = 0, sin(90) = 1
	// m[0][0] = 0, m[1][0] = -1
	// m[0][1] = 1, m[1][1] = 0
	m := NewRotationMatrix(ROTZ, math.Pi/2)

	if math.Abs(m.ThisMatrix[0][0]) > 1e-9 {
		t.Errorf("ROTZ m[0][0] = %f; want 0", m.ThisMatrix[0][0])
	}
	if math.Abs(m.ThisMatrix[1][0]-(-1.0)) > 1e-9 {
		t.Errorf("ROTZ m[1][0] = %f; want -1", m.ThisMatrix[1][0])
	}
}

func TestIdentMatrix(t *testing.T) {
	m := IdentMatrix()
	for i := 0; i < 4; i++ {
		if m.ThisMatrix[i][i] != 1.0 {
			t.Errorf("IdentMatrix()[%d][%d] = %f; want 1.0", i, i, m.ThisMatrix[i][i])
		}
	}
}

func TestTransMatrix(t *testing.T) {
	m := TransMatrix(1, 2, 3)
	if m.ThisMatrix[3][0] != 1 || m.ThisMatrix[3][1] != 2 || m.ThisMatrix[3][2] != 3 {
		t.Errorf("TransMatrix(1, 2, 3) translation part mismatch: %v", m.ThisMatrix[3])
	}
	if m.ThisMatrix[0][0] != 1 {
		t.Error("TransMatrix() should be identity-like")
	}
}

// TestMatrix_AddRow is removed.

func TestMatrix_MultiplyBy(t *testing.T) {
	// Identity * Identity = Identity
	m1 := IdentMatrix()
	m2 := IdentMatrix()
	m3 := m1.MultiplyBy(m2)

	for i := 0; i < 4; i++ {
		if m3.ThisMatrix[i][i] != 1.0 {
			t.Errorf("Identity * Identity result[%d][%d] = %f; want 1.0", i, i, m3.ThisMatrix[i][i])
		}
	}
}

func TestMatrix_TransformObj(t *testing.T) {
	// Translate point (1, 1, 1) by (1, 2, 3) -> (2, 3, 4)
	trans := TransMatrix(1, 2, 3)
	src := []Vector3{{X: 1, Y: 1, Z: 1, W: 1.0}}
	dest := []Vector3{{X: 0, Y: 0, Z: 0, W: 0.0}}

	trans.TransformObj(src, dest)

	if dest[0].X != 2 || dest[0].Y != 3 || dest[0].Z != 4 {
		t.Errorf("TransformObj() resulted in %v; want {2, 3, 4}", dest[0])
	}
}

func TestMatrix_Copy(t *testing.T) {
	m := IdentMatrix()
	c := m.Copy()

	c.ThisMatrix[0][0] = 5
	if m.ThisMatrix[0][0] != 1 {
		t.Error("Copy() did not deep copy")
	}
}

func TestMatrix_RotateVector3(t *testing.T) {
	// Rotate (1, 0, 0) 90 deg around Z -> (0, 1, 0)
	m := NewRotationMatrix(ROTZ, math.Pi/2)
	v := NewVector3(1, 0, 0)
	rotated := m.RotateVector3(v)

	if math.Abs(rotated.X) > 1e-9 {
		t.Errorf("RotateVector3 X = %f; want 0", rotated.X)
	}
	if math.Abs(rotated.Y-1.0) > 1e-9 {
		t.Errorf("RotateVector3 Y = %f; want 1", rotated.Y)
	}
}

func matricesEqual(m1, m2 Matrix) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(m1.ThisMatrix[i][j]-m2.ThisMatrix[i][j]) > 1e-6 {
				return false
			}
		}
	}
	return true
}

func TestNewRotationMatrixAxis(t *testing.T) {
	// Case 1: Rotate around X by Pi/2
	axisX := NewVector3(1, 0, 0)
	m1 := NewRotationMatrixAxis(axisX, math.Pi/2)
	expected1 := NewRotationMatrix(ROTX, math.Pi/2)

	if !matricesEqual(m1, expected1) {
		t.Errorf("Rotation around X axis mismatch.\nGot:\n%v\nExpected:\n%v", m1, expected1)
	}

	// Case 2: Rotate around Y by Pi/2
	axisY := NewVector3(0, 1, 0)
	m2 := NewRotationMatrixAxis(axisY, math.Pi/2)
	expected2 := NewRotationMatrix(ROTY, math.Pi/2)

	if !matricesEqual(m2, expected2) {
		t.Errorf("Rotation around Y axis mismatch.\nGot:\n%v\nExpected:\n%v", m2, expected2)
	}

	// Case 3: Rotate vector (1, 0, 0) around (1, 1, 1) by 2*Pi/3 -> (0, 1, 0)
	// (1,1,1) normalized is (1/sqrt3, 1/sqrt3, 1/sqrt3).
	// 2*Pi/3 is 120 degrees.
	// This rotation permutes axes X->Y->Z->X.
	// So (1,0,0) -> (0,1,0).
	axisDiag := NewVector3(1, 1, 1)
	m3 := NewRotationMatrixAxis(axisDiag, 2*math.Pi/3)
	v := NewVector3(1, 0, 0)
	rotated := m3.RotateVector3(v)
	expectedV := NewVector3(0, 1, 0)

	if math.Abs(rotated.X-expectedV.X) > 1e-6 ||
		math.Abs(rotated.Y-expectedV.Y) > 1e-6 ||
		math.Abs(rotated.Z-expectedV.Z) > 1e-6 {
		t.Errorf("Rotation around diagonal failed. Got %v, expected %v", rotated, expectedV)
	}
}
