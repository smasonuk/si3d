package si3d

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	camMatrixRev   Matrix
	cameraPosition Vector3
	cameraRotation mgl64.Quat
	NearPlane      float64
}

func (c *Camera) GetNearPlane() float64 {
	return c.NearPlane
}

func NewCamera(xp, yp, zp, xa, ya, za float64) *Camera {
	c := &Camera{}
	c.NearPlane = 10.0

	// NEW Logic
	rotX := mgl64.QuatRotate(xa, mgl64.Vec3{1, 0, 0})
	rotY := mgl64.QuatRotate(ya, mgl64.Vec3{0, 1, 0})
	rotZ := mgl64.QuatRotate(za, mgl64.Vec3{0, 0, 1})
	c.cameraRotation = rotX.Mul(rotY).Mul(rotZ)

	c.cameraPosition = NewVector3(xp, yp, zp)

	c.updateMatrix()

	return c
}

func (c *Camera) updateMatrix() {
	sTransWorldToCamera := TransMatrix(-c.cameraPosition.X, -c.cameraPosition.Y, -c.cameraPosition.Z)
	invRot := c.cameraRotation.Conjugate()
	rotMat := ToGoSieMatrixFromQuat(invRot)
	c.camMatrixRev = rotMat.MultiplyBy(sTransWorldToCamera)
}

func (c *Camera) GetCameraMatrix() Matrix {
	return c.camMatrixRev
}

func NewCameraLookAt(camPos Vector3, lookAt Vector3, up Vector3) *Camera {
	lookAtMat := mgl64.LookAt(
		lookAt.X, lookAt.Y, lookAt.Z,
		camPos.X, camPos.Y, camPos.Z,

		0, 1, 0,
	)

	sMat := ToGoSieMatrix(lookAtMat)

	c := &Camera{
		camMatrixRev:   sMat,
		cameraPosition: NewVector3(camPos.X, camPos.Y, camPos.Z),
		cameraRotation: mgl64.QuatIdent(),
		NearPlane:      10.0,
	}

	return c
}

func NewCameraLookMatrixAt3(cameraLocation Vector3, lookAt Vector3, up Vector3) Matrix {
	sTransWorldToCamera := TransMatrix(
		-cameraLocation.X,
		-cameraLocation.Y,
		-cameraLocation.Z,
	)

	dirY := lookAt.Y - cameraLocation.Y
	dirX := lookAt.X - cameraLocation.X
	dirZ := lookAt.Z - cameraLocation.Z
	lookADirVec := NewVector3(dirX, dirY, dirZ)

	angleRadY := angleY(lookAt, cameraLocation)
	sMatY := NewRotationMatrix(ROTY, angleRadY)
	newLookAtDirVec := sMatY.RotateVector3(lookADirVec)

	angleDown := angleDown(newLookAtDirVec)

	sMatX := NewRotationMatrix(ROTX, angleDown)

	sMat := sMatX.MultiplyBy(sMatY.MultiplyBy(sTransWorldToCamera))

	return sMat
}

func angleDown(lookADirVec Vector3) float64 {
	hypot := math.Sqrt(lookADirVec.X*lookADirVec.X + lookADirVec.Y*lookADirVec.Y + lookADirVec.Z*lookADirVec.Z)
	adjacent := lookADirVec.Y

	angleDownRad := math.Acos(adjacent / hypot)        // Angle in radians
	angleDownDegrees := angleDownRad * (180 / math.Pi) // Convert to degrees
	angleDownDegrees = 90 - angleDownDegrees

	return degreesToRadians(angleDownDegrees)

}

func angleY(lookAt Vector3, cameraLocation Vector3) float64 {
	dirY := lookAt.Z - cameraLocation.Z
	dirX := lookAt.X - cameraLocation.X
	angleY := math.Atan2(dirY, dirX)
	angleY = angleY - math.Pi/2 // Adjust to match the camera's forward direction
	return angleY
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func (c *Camera) LookAt(lookAt Vector3, up Vector3) {

	sMat := NewCameraLookMatrixAt3(c.cameraPosition, lookAt, up)

	c.camMatrixRev = sMat
}

func (c *Camera) GetPosition() Vector3 {
	return c.cameraPosition
}

func (c *Camera) SetCameraPosition(x, y, z float64) {
	c.cameraPosition = NewVector3(x, y, z)
}

func (c *Camera) AddXPosition(x float64) {
	c.cameraPosition.X += x
}

func (c *Camera) AddYPosition(y float64) {
	c.cameraPosition.Y += y
}

func (c *Camera) AddZPosition(z float64) {
	c.cameraPosition.Z += z
}

func (c *Camera) GetMatrix() Matrix {
	return c.camMatrixRev
}

func (c *Camera) AddAngle(x, y, z float64) {
	rotX := mgl64.QuatRotate(x, mgl64.Vec3{1, 0, 0})
	rotY := mgl64.QuatRotate(y, mgl64.Vec3{0, 1, 0})
	rotZ := mgl64.QuatRotate(z, mgl64.Vec3{0, 0, 1})

	deltaRot := rotX.Mul(rotY).Mul(rotZ)

	c.cameraRotation = c.cameraRotation.Mul(deltaRot)

	c.updateMatrix()
}
