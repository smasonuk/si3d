package gosie3d

import "github.com/go-gl/mathgl/mgl64"

type Transform struct {
	Position Vector3
	Rotation mgl64.Quat
	Scale    Vector3
}

func NewTransform() *Transform {
	return &Transform{
		Scale:    NewVector3(1.0, 1.0, 1.0),
		Rotation: mgl64.QuatIdent(),
	}
}

// GetMatrix returns the combined local matrix:
// M = Scale * Rotation * Trans
// This corresponds to applying Scale, then Rotation, then Translation to a row vector.
func (t *Transform) GetMatrix() Matrix {
	// Translation matrix
	m := TransMatrix(t.Position.X, t.Position.Y, t.Position.Z)

	// Multiply by Rotations
	// Note: MultiplyBy(A) computes A * m (where m is the receiver)
	// We want Scale * Rot * Trans.
	// We start with Trans (T).
	// T.MultiplyBy(Rot) -> Rot * T.
	// (Rot * T).MultiplyBy(Scale) -> Scale * Rot * T.

	rotMat := ToGoSieMatrixFromQuat(t.Rotation)
	m = m.MultiplyBy(rotMat)
	m = m.MultiplyBy(ScaleMatrix(t.Scale.X, t.Scale.Y, t.Scale.Z))

	return m
}

func (t *Transform) Rotate(axis Vector3, angleRadians float64) {
	// Global rotation: Q_new = Rotation * Q_rot
	rot := mgl64.QuatRotate(angleRadians, mgl64.Vec3{axis.X, axis.Y, axis.Z})
	t.Rotation = t.Rotation.Mul(rot)
}

func (t *Transform) RotateLocal(axis Vector3, angleRadians float64) {
	// Local rotation: Q_new = Q_rot * Rotation
	rot := mgl64.QuatRotate(angleRadians, mgl64.Vec3{axis.X, axis.Y, axis.Z})
	t.Rotation = rot.Mul(t.Rotation)
}
