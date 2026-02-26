package si3d

import (
	"math"
)

func GetLength(vec []float64) float64 {
	return math.Sqrt(vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2])
}
