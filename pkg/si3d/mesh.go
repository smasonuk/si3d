package si3d

type Mesh struct {
	Points     []Vector3
	pointIndex map[Vector3]int
}

func NewMesh() *Mesh {
	return &Mesh{
		Points:     make([]Vector3, 0, 100), // pre-allocate some space
		pointIndex: make(map[Vector3]int),
	}

}

// AddPoint now uses the map for an average O(1) lookup.
func (m *Mesh) AddPoint(v Vector3) (Vector3, int) {
	// Check if the point already exists using the map.
	// Vector3 is comparable (contains only float64).
	// Note: W component is included in comparison.
	if index, found := m.pointIndex[v]; found {
		return m.Points[index], index
	}

	// If not found, add the new point.
	m.Points = append(m.Points, v)
	newIndex := len(m.Points) - 1

	// Add the new point's key and index to the map.
	m.pointIndex[v] = newIndex

	return v, newIndex
}

// Copy must also duplicate the pointIndex map.
func (m *Mesh) Copy() *Mesh {
	newPointIndex := make(map[Vector3]int, len(m.pointIndex))
	for key, value := range m.pointIndex {
		newPointIndex[key] = value
	}

	newPoints := make([]Vector3, len(m.Points))
	copy(newPoints, m.Points)

	return &Mesh{
		Points:     newPoints,
		pointIndex: newPointIndex,
	}
}
