package gosie3d

import (
	"testing"
)

func TestNewMesh(t *testing.T) {
	m := NewMesh()
	if m == nil {
		t.Fatal("NewMesh returned nil")
	}
	if len(m.Points) != 0 {
		t.Errorf("NewMesh should have 0 points, got %d", len(m.Points))
	}
	if m.pointIndex == nil {
		t.Errorf("NewMesh pointIndex map should be initialized")
	}
}

func TestMesh_AddPoint(t *testing.T) {
	m := NewMesh()
	p1 := NewVector3(1, 2, 3)

	// Add first point
	v1, idx1 := m.AddPoint(p1)
	if idx1 != 0 {
		t.Errorf("First point index should be 0, got %d", idx1)
	}
	if v1.X != 1 || v1.Y != 2 || v1.Z != 3 {
		t.Errorf("AddPoint returned incorrect vector: %v", v1)
	}

	// Add same point again
	v2, idx2 := m.AddPoint(p1)
	if idx2 != 0 {
		t.Errorf("Duplicate point should return index 0, got %d", idx2)
	}
	if v2.X != 1 || v2.Y != 2 || v2.Z != 3 {
		t.Errorf("Duplicate point returned incorrect vector: %v", v2)
	}
	if len(m.Points) != 1 {
		t.Errorf("Duplicate point should not increase points count, got %d", len(m.Points))
	}

	// Add different point
	p2 := NewVector3(4, 5, 6)
	v3, idx3 := m.AddPoint(p2)
	if idx3 != 1 {
		t.Errorf("Second point index should be 1, got %d", idx3)
	}
	if v3.X != 4 || v3.Y != 5 || v3.Z != 6 {
		t.Errorf("Second point returned incorrect vector: %v", v3)
	}
	if len(m.Points) != 2 {
		t.Errorf("Points count should be 2, got %d", len(m.Points))
	}
}

func TestMesh_Copy(t *testing.T) {
	m := NewMesh()
	p1 := NewVector3(1, 2, 3)
	m.AddPoint(p1)

	mCopy := m.Copy()
	if mCopy == nil {
		t.Fatal("Copy returned nil")
	}
	if len(mCopy.Points) != len(m.Points) {
		t.Errorf("Copy Points length mismatch")
	}
	if len(mCopy.pointIndex) != len(m.pointIndex) {
		t.Errorf("Copy pointIndex map length mismatch")
	}

	// Check pointIndex map content
	idx, found := mCopy.pointIndex[p1]
	if !found {
		t.Errorf("Copy pointIndex map missing entry")
	}
	if idx != 0 {
		t.Errorf("Copy pointIndex map incorrect index")
	}

	// Verify deep copy of slices/maps
	// Modifying original shouldn't affect copy
	p2 := NewVector3(4, 5, 6)
	m.AddPoint(p2)
	if len(mCopy.Points) == len(m.Points) {
		t.Errorf("Copy should be deep copy for Points slice")
	}
}
