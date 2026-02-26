package si3d

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

func NewMatrix() Matrix {
	// Returns a zero-initialized matrix (all elements 0.0)
	return Matrix{}
}

type Matrix struct {
	// Fixed 4x4 array for better performance and cache locality.
	ThisMatrix [4][4]float64
}

const (
	ROTX = 0
	ROTY = 1
	ROTZ = 2
)

func NewMatrixFromData(aMatrix [][]float64) Matrix {
	m := Matrix{}
	for i := range aMatrix {
		if i >= 4 {
			break
		}
		for j := range aMatrix[i] {
			if j >= 4 {
				break
			}
			m.ThisMatrix[i][j] = aMatrix[i][j]
		}
	}
	return m
}

func NewRotationMatrix(aRotation int, theta float64) Matrix {
	m := Matrix{}
	switch aRotation {
	case ROTX:
		m.ThisMatrix[0][0] = 1.0
		m.ThisMatrix[1][1] = math.Cos(theta)
		m.ThisMatrix[2][1] = -math.Sin(theta)
		m.ThisMatrix[1][2] = math.Sin(theta)
		m.ThisMatrix[2][2] = math.Cos(theta)
		m.ThisMatrix[3][3] = 1.0
	case ROTY:
		m.ThisMatrix[0][0] = math.Cos(theta)
		m.ThisMatrix[2][0] = math.Sin(theta)
		m.ThisMatrix[0][2] = -math.Sin(theta)
		m.ThisMatrix[2][2] = math.Cos(theta)
		m.ThisMatrix[1][1] = 1.0
		m.ThisMatrix[3][3] = 1.0
	case ROTZ:
		m.ThisMatrix[2][2] = 1.0
		m.ThisMatrix[3][3] = 1.0
		m.ThisMatrix[0][0] = math.Cos(theta)
		m.ThisMatrix[1][0] = -math.Sin(theta)
		m.ThisMatrix[0][1] = math.Sin(theta)
		m.ThisMatrix[1][1] = math.Cos(theta)
	}
	return m
}

func NewRotationMatrixAxis(axis Vector3, angleRadians float64) Matrix {
	axis = axis.Normalize()
	kx, ky, kz := axis.X, axis.Y, axis.Z
	c := math.Cos(angleRadians)
	s := math.Sin(angleRadians)
	omc := 1.0 - c

	m := Matrix{}
	// Note: The signs of the sine terms are flipped compared to the standard
	// column-vector Rodrigues formula to account for row-vector multiplication (v * M).
	m.ThisMatrix[0][0] = c + kx*kx*omc
	m.ThisMatrix[0][1] = kx*ky*omc + kz*s
	m.ThisMatrix[0][2] = kx*kz*omc - ky*s
	m.ThisMatrix[0][3] = 0.0

	m.ThisMatrix[1][0] = ky*kx*omc - kz*s
	m.ThisMatrix[1][1] = c + ky*ky*omc
	m.ThisMatrix[1][2] = ky*kz*omc + kx*s
	m.ThisMatrix[1][3] = 0.0

	m.ThisMatrix[2][0] = kz*kx*omc + ky*s
	m.ThisMatrix[2][1] = kz*ky*omc - kx*s
	m.ThisMatrix[2][2] = c + kz*kz*omc
	m.ThisMatrix[2][3] = 0.0

	m.ThisMatrix[3][0] = 0.0
	m.ThisMatrix[3][1] = 0.0
	m.ThisMatrix[3][2] = 0.0
	m.ThisMatrix[3][3] = 1.0
	return m
}

func IdentMatrix() Matrix {
	m := Matrix{}
	m.ThisMatrix[0][0] = 1.0
	m.ThisMatrix[1][1] = 1.0
	m.ThisMatrix[2][2] = 1.0
	m.ThisMatrix[3][3] = 1.0
	return m
}

func ScaleMatrix(x, y, z float64) Matrix {
	m := Matrix{}
	m.ThisMatrix[0][0] = x
	m.ThisMatrix[1][1] = y
	m.ThisMatrix[2][2] = z
	m.ThisMatrix[3][3] = 1.0
	return m
}

func TransMatrix(x, y, z float64) Matrix {
	nm := Matrix{}
	nm.ThisMatrix[3][0] = x
	nm.ThisMatrix[3][1] = y
	nm.ThisMatrix[3][2] = z
	nm.ThisMatrix[0][0] = 1.0
	nm.ThisMatrix[1][1] = 1.0
	nm.ThisMatrix[2][2] = 1.0
	nm.ThisMatrix[3][3] = 1.0
	return nm
}

// AddRow is removed as Matrix is now fixed size.

func (m Matrix) MultiplyBy(aMatrix Matrix) Matrix {
	newMatrix := Matrix{}

	newMatrix.ThisMatrix[0][0] = m.ThisMatrix[0][0]*aMatrix.ThisMatrix[0][0] + m.ThisMatrix[1][0]*aMatrix.ThisMatrix[0][1] + m.ThisMatrix[2][0]*aMatrix.ThisMatrix[0][2] + m.ThisMatrix[3][0]*aMatrix.ThisMatrix[0][3]
	newMatrix.ThisMatrix[1][0] = m.ThisMatrix[0][0]*aMatrix.ThisMatrix[1][0] + m.ThisMatrix[1][0]*aMatrix.ThisMatrix[1][1] + m.ThisMatrix[2][0]*aMatrix.ThisMatrix[1][2] + m.ThisMatrix[3][0]*aMatrix.ThisMatrix[1][3]
	newMatrix.ThisMatrix[2][0] = m.ThisMatrix[0][0]*aMatrix.ThisMatrix[2][0] + m.ThisMatrix[1][0]*aMatrix.ThisMatrix[2][1] + m.ThisMatrix[2][0]*aMatrix.ThisMatrix[2][2] + m.ThisMatrix[3][0]*aMatrix.ThisMatrix[2][3]
	newMatrix.ThisMatrix[3][0] = m.ThisMatrix[0][0]*aMatrix.ThisMatrix[3][0] + m.ThisMatrix[1][0]*aMatrix.ThisMatrix[3][1] + m.ThisMatrix[2][0]*aMatrix.ThisMatrix[3][2] + m.ThisMatrix[3][0]*aMatrix.ThisMatrix[3][3]
	newMatrix.ThisMatrix[0][1] = m.ThisMatrix[0][1]*aMatrix.ThisMatrix[0][0] + m.ThisMatrix[1][1]*aMatrix.ThisMatrix[0][1] + m.ThisMatrix[2][1]*aMatrix.ThisMatrix[0][2] + m.ThisMatrix[3][1]*aMatrix.ThisMatrix[0][3]
	newMatrix.ThisMatrix[1][1] = m.ThisMatrix[0][1]*aMatrix.ThisMatrix[1][0] + m.ThisMatrix[1][1]*aMatrix.ThisMatrix[1][1] + m.ThisMatrix[2][1]*aMatrix.ThisMatrix[1][2] + m.ThisMatrix[3][1]*aMatrix.ThisMatrix[1][3]
	newMatrix.ThisMatrix[2][1] = m.ThisMatrix[0][1]*aMatrix.ThisMatrix[2][0] + m.ThisMatrix[1][1]*aMatrix.ThisMatrix[2][1] + m.ThisMatrix[2][1]*aMatrix.ThisMatrix[2][2] + m.ThisMatrix[3][1]*aMatrix.ThisMatrix[2][3]
	newMatrix.ThisMatrix[3][1] = m.ThisMatrix[0][1]*aMatrix.ThisMatrix[3][0] + m.ThisMatrix[1][1]*aMatrix.ThisMatrix[3][1] + m.ThisMatrix[2][1]*aMatrix.ThisMatrix[3][2] + m.ThisMatrix[3][1]*aMatrix.ThisMatrix[3][3]
	newMatrix.ThisMatrix[0][2] = m.ThisMatrix[0][2]*aMatrix.ThisMatrix[0][0] + m.ThisMatrix[1][2]*aMatrix.ThisMatrix[0][1] + m.ThisMatrix[2][2]*aMatrix.ThisMatrix[0][2] + m.ThisMatrix[3][2]*aMatrix.ThisMatrix[0][3]
	newMatrix.ThisMatrix[1][2] = m.ThisMatrix[0][2]*aMatrix.ThisMatrix[1][0] + m.ThisMatrix[1][2]*aMatrix.ThisMatrix[1][1] + m.ThisMatrix[2][2]*aMatrix.ThisMatrix[1][2] + m.ThisMatrix[3][2]*aMatrix.ThisMatrix[1][3]
	newMatrix.ThisMatrix[2][2] = m.ThisMatrix[0][2]*aMatrix.ThisMatrix[2][0] + m.ThisMatrix[1][2]*aMatrix.ThisMatrix[2][1] + m.ThisMatrix[2][2]*aMatrix.ThisMatrix[2][2] + m.ThisMatrix[3][2]*aMatrix.ThisMatrix[2][3]
	newMatrix.ThisMatrix[3][2] = m.ThisMatrix[0][2]*aMatrix.ThisMatrix[3][0] + m.ThisMatrix[1][2]*aMatrix.ThisMatrix[3][1] + m.ThisMatrix[2][2]*aMatrix.ThisMatrix[3][2] + m.ThisMatrix[3][2]*aMatrix.ThisMatrix[3][3]
	newMatrix.ThisMatrix[0][3] = m.ThisMatrix[0][3]*aMatrix.ThisMatrix[0][0] + m.ThisMatrix[1][3]*aMatrix.ThisMatrix[0][1] + m.ThisMatrix[2][3]*aMatrix.ThisMatrix[0][2] + m.ThisMatrix[3][3]*aMatrix.ThisMatrix[0][3]
	newMatrix.ThisMatrix[1][3] = m.ThisMatrix[0][3]*aMatrix.ThisMatrix[1][0] + m.ThisMatrix[1][3]*aMatrix.ThisMatrix[1][1] + m.ThisMatrix[2][3]*aMatrix.ThisMatrix[1][2] + m.ThisMatrix[3][3]*aMatrix.ThisMatrix[1][3]
	newMatrix.ThisMatrix[2][3] = m.ThisMatrix[0][3]*aMatrix.ThisMatrix[2][0] + m.ThisMatrix[1][3]*aMatrix.ThisMatrix[2][1] + m.ThisMatrix[2][3]*aMatrix.ThisMatrix[2][2] + m.ThisMatrix[3][3]*aMatrix.ThisMatrix[2][3]
	newMatrix.ThisMatrix[3][3] = m.ThisMatrix[0][3]*aMatrix.ThisMatrix[3][0] + m.ThisMatrix[1][3]*aMatrix.ThisMatrix[3][1] + m.ThisMatrix[2][3]*aMatrix.ThisMatrix[3][2] + m.ThisMatrix[3][3]*aMatrix.ThisMatrix[3][3]
	return newMatrix
}

func (m *Matrix) TransformObj(src, dest []Vector3) {
	numItems := len(src)
	if numItems == 0 {
		return
	}

	if numItems < 1000 {
		for x := 0; x < numItems; x++ {
			sx, sy, sz := src[x].X, src[x].Y, src[x].Z
			dest[x].X = m.ThisMatrix[0][0]*sx + m.ThisMatrix[1][0]*sy + m.ThisMatrix[2][0]*sz + m.ThisMatrix[3][0]
			dest[x].Y = m.ThisMatrix[0][1]*sx + m.ThisMatrix[1][1]*sy + m.ThisMatrix[2][1]*sz + m.ThisMatrix[3][1]
			dest[x].Z = m.ThisMatrix[0][2]*sx + m.ThisMatrix[1][2]*sy + m.ThisMatrix[2][2]*sz + m.ThisMatrix[3][2]
			dest[x].W = 1.0
		}
		return
	}

	numCPU := runtime.NumCPU()
	chunkSize := (numItems + numCPU - 1) / numCPU

	var wg sync.WaitGroup

	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > numItems {
			end = numItems
		}
		if start >= end {
			continue
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for x := start; x < end; x++ {
				sx, sy, sz := src[x].X, src[x].Y, src[x].Z
				dest[x].X = m.ThisMatrix[0][0]*sx + m.ThisMatrix[1][0]*sy + m.ThisMatrix[2][0]*sz + m.ThisMatrix[3][0]
				dest[x].Y = m.ThisMatrix[0][1]*sx + m.ThisMatrix[1][1]*sy + m.ThisMatrix[2][1]*sz + m.ThisMatrix[3][1]
				dest[x].Z = m.ThisMatrix[0][2]*sx + m.ThisMatrix[1][2]*sy + m.ThisMatrix[2][2]*sz + m.ThisMatrix[3][2]
				dest[x].W = 1.0
			}
		}(start, end)
	}

	wg.Wait()
}

func (m Matrix) Copy() Matrix {
	newMat := Matrix{}
	newMat.ThisMatrix = m.ThisMatrix // Array copy is value copy
	return newMat
}

func (m *Matrix) String() string {
	var sb strings.Builder
	for i := 0; i < 4; i++ {
		if i > 0 {
			sb.WriteString("\n")
		}
		for j := 0; j < 4; j++ {
			if j > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%f", m.ThisMatrix[i][j]))
		}
	}
	return sb.String()
}

func ToGoSieMatrix(m mgl64.Mat4) Matrix {
	return NewMatrixFromData(
		[][]float64{
			{m[0], m[1], m[2], m[3]},
			{m[4], m[5], m[6], m[7]},
			{m[8], m[9], m[10], m[11]},
			{m[12], m[13], m[14], m[15]},
		},
	)
}

func ToGoSieMatrixFromQuat(q mgl64.Quat) Matrix {
	return ToGoSieMatrix(q.Mat4())
}

// RotateVector3 rotates a Vector3 by the matrix's 3x3 rotation component.
// It does not apply translation, making it suitable for direction vectors.
func (m Matrix) RotateVector3(v Vector3) Vector3 {
	// Extract the source vector components for clarity.
	vx, vy, vz := v.X, v.Y, v.Z

	// Apply the 3x3 rotation part of the matrix. This is a standard
	// vector-matrix multiplication that ignores the translation part of the matrix.
	newX := m.ThisMatrix[0][0]*vx + m.ThisMatrix[1][0]*vy + m.ThisMatrix[2][0]*vz
	newY := m.ThisMatrix[0][1]*vx + m.ThisMatrix[1][1]*vy + m.ThisMatrix[2][1]*vz
	newZ := m.ThisMatrix[0][2]*vx + m.ThisMatrix[1][2]*vy + m.ThisMatrix[2][2]*vz

	// Return a new Vector3 with the rotated coordinates.
	return NewVector3(newX, newY, newZ)
}

func (m *Matrix) TransformNormals(src, dest []Vector3) {
	numItems := len(src)
	if numItems == 0 {
		return
	}

	if numItems < 1000 {
		for x := 0; x < numItems; x++ {
			sx, sy, sz := src[x].X, src[x].Y, src[x].Z

			// This is a 3x3 rotation of a vector, it deliberately ignores the
			// translation components of the matrix (m.ThisMatrix[3][...]).
			dest[x].X = m.ThisMatrix[0][0]*sx + m.ThisMatrix[1][0]*sy + m.ThisMatrix[2][0]*sz
			dest[x].Y = m.ThisMatrix[0][1]*sx + m.ThisMatrix[1][1]*sy + m.ThisMatrix[2][1]*sz
			dest[x].Z = m.ThisMatrix[0][2]*sx + m.ThisMatrix[1][2]*sy + m.ThisMatrix[2][2]*sz
		}
		return
	}

	numCPU := runtime.NumCPU()
	chunkSize := (numItems + numCPU - 1) / numCPU

	var wg sync.WaitGroup

	for i := 0; i < numCPU; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > numItems {
			end = numItems
		}
		if start >= end {
			continue
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for x := start; x < end; x++ {
				sx, sy, sz := src[x].X, src[x].Y, src[x].Z

				// This is a 3x3 rotation of a vector, it deliberately ignores the
				// translation components of the matrix (m.ThisMatrix[3][...]).
				dest[x].X = m.ThisMatrix[0][0]*sx + m.ThisMatrix[1][0]*sy + m.ThisMatrix[2][0]*sz
				dest[x].Y = m.ThisMatrix[0][1]*sx + m.ThisMatrix[1][1]*sy + m.ThisMatrix[2][1]*sz
				dest[x].Z = m.ThisMatrix[0][2]*sx + m.ThisMatrix[1][2]*sy + m.ThisMatrix[2][2]*sz
			}
		}(start, end)
	}

	wg.Wait()
}
