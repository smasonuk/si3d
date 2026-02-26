package gosie3d

import "math"

type Vector3 struct {
	X float64
	Y float64
	Z float64
	W float64
}

// Add returns a new Vector3 that is the sum of v and other.
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
		W: 1.0,
	}
}

func NewVector3(x, y, z float64) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
		W: 1.0, // W component is typically 1.0 for vectors
	}
}

func NewVector3Full(x, y, z, w float64) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}

func NewVector3dFromArray(normal []float64) Vector3 {
	return Vector3{
		X: normal[0],
		Y: normal[1],
		Z: normal[2],
		W: 1.0,
	}
}

func (v Vector3) Normalize() Vector3 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if length == 0 {
		return v
	}
	return Vector3{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
		W: v.W,
	}
}

func (v Vector3) Copy() Vector3 {
	return v
}

// DistanceTo
func (v Vector3) DistanceTo(other Vector3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// DistanceSquaredTo
func (v Vector3) DistanceSquaredTo(other Vector3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

// Subtract returns a new Vector3 that is the difference of v1 and v2.
func Subtract(v1, v2 Vector3) Vector3 {
	return NewVector3(
		v1.X-v2.X,
		v1.Y-v2.Y,
		v1.Z-v2.Z,
	)
}

// Cross computes the Cross product of two vectors and returns a new Vector3.
func Cross(v1, v2 Vector3) Vector3 {
	return NewVector3(
		v1.Y*v2.Z-v1.Z*v2.Y,
		v1.Z*v2.X-v1.X*v2.Z,
		v1.X*v2.Y-v1.Y*v2.X,
	)
}

// Dot computes the dot product of two vectors.
func Dot(v1, v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}
