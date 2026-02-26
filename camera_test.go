package gosie3d

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestNewCamera(t *testing.T) {
	c := NewCamera(10, 20, 30, 0, 0, 0)
	if c == nil {
		t.Fatal("NewCamera returned nil")
	}

	pos := c.GetPosition()
	if pos.X != 10 || pos.Y != 20 || pos.Z != 30 {
		t.Errorf("NewCamera position expected (10, 20, 30), got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}
}

func TestCamera_SetCameraPosition(t *testing.T) {
	c := NewCamera(0, 0, 0, 0, 0, 0)
	c.SetCameraPosition(100, 200, 300)

	pos := c.GetPosition()
	if pos.X != 100 || pos.Y != 200 || pos.Z != 300 {
		t.Errorf("SetCameraPosition failed, got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}
}

func TestCamera_AddPosition(t *testing.T) {
	c := NewCamera(10, 20, 30, 0, 0, 0)
	c.AddXPosition(5)
	c.AddYPosition(-5)
	c.AddZPosition(10)

	pos := c.GetPosition()
	if pos.X != 15 || pos.Y != 15 || pos.Z != 40 {
		t.Errorf("AddPosition failed, got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}
}

func TestNewCameraLookAt(t *testing.T) {
	camPos := NewVector3(0, 0, 0)
	lookAt := NewVector3(0, 0, -10)
	up := NewVector3(0, 1, 0)

	c := NewCameraLookAt(camPos, lookAt, up)
	if c == nil {
		t.Fatal("NewCameraLookAt returned nil")
	}

	pos := c.GetPosition()
	if pos.X != 0 || pos.Y != 0 || pos.Z != 0 {
		t.Errorf("NewCameraLookAt position incorrect, got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}

	// Verify lookAt matrix logic indirectly via basic properties or just ensuring it's not nil
	c.GetCameraMatrix()
}

func TestHelper_degreesToRadians(t *testing.T) {
	rad := degreesToRadians(180)
	if math.Abs(rad-math.Pi) > 1e-6 {
		t.Errorf("degreesToRadians(180) = %f, expected %f", rad, math.Pi)
	}
}

func TestHelper_angleY(t *testing.T) {
	// Test angleY logic
	// angleY calculates angle between two points in XZ plane
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(1, 0, 1) // 45 degrees

	// angleY implementation:
	// dirY := lookAt.Z - cameraLocation.Z
	// dirX := lookAt.X - cameraLocation.X
	// angleY := math.Atan2(dirY, dirX)
	// angleY = angleY - math.Pi/2

	// For (1,0,1) from (0,0,0):
	// dirY = 1, dirX = 1 => Atan2(1,1) = Pi/4
	// result = Pi/4 - Pi/2 = -Pi/4

	ang := angleY(p2, p1)
	expected := math.Pi/4 - math.Pi/2
	if math.Abs(ang-expected) > 1e-6 {
		t.Errorf("angleY expected %f, got %f", expected, ang)
	}
}

func TestCamera_AddAngle(t *testing.T) {
	c := NewCamera(0, 0, 0, 0, 0, 0)
	c.AddAngle(0.1, 0.2, 0.3)

	// Since cameraAngle was replaced by cameraRotation (Quat), we verify the quaternion.
	// Initial rotation is Identity.
	// AddAngle applies Qx(0.1) * Qy(0.2) * Qz(0.3).

	rotX := mgl64.QuatRotate(0.1, mgl64.Vec3{1, 0, 0})
	rotY := mgl64.QuatRotate(0.2, mgl64.Vec3{0, 1, 0})
	rotZ := mgl64.QuatRotate(0.3, mgl64.Vec3{0, 0, 1})

	// c.cameraRotation was Identity.
	// AddAngle -> c.cameraRotation = c.cameraRotation.Mul(dQ)
	expected := rotX.Mul(rotY).Mul(rotZ)

	q := c.cameraRotation

	if math.Abs(q.W-expected.W) > 1e-6 ||
		math.Abs(q.V.X()-expected.V.X()) > 1e-6 ||
		math.Abs(q.V.Y()-expected.V.Y()) > 1e-6 ||
		math.Abs(q.V.Z()-expected.V.Z()) > 1e-6 {
		t.Errorf("AddAngle failed. Got %v, expected %v", q, expected)
	}

	// Check if matrix updated (smoke test)
	c.GetMatrix()
}
