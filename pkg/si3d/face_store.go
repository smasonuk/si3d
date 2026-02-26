package si3d

import "sort"

type FaceStore struct {
	faces []*Face
}

func NewFaceStore() *FaceStore {
	return &FaceStore{faces: make([]*Face, 0, 10)}
}
func (fs *FaceStore) AddFace(f *Face) {
	fs.faces = append(fs.faces, f)
}
func (fs *FaceStore) GetFace(i int) *Face {
	return fs.faces[i]
}
func (fs *FaceStore) FaceCount() int {
	return len(fs.faces)
}
func (fs *FaceStore) RemoveFaceAt(i int) *Face {
	if i < 0 || i >= len(fs.faces) {
		return nil
	}
	f := fs.faces[i]
	fs.faces = append(fs.faces[:i], fs.faces[i+1:]...)
	return f
}

// sort the faces so that the faces farther away are at the start of the slice
func (fs *FaceStore) SortFacesByDistance(pos Vector3) {
	if len(fs.faces) == 0 {
		return
	}

	type sortableFace struct {
		f      *Face
		distSq float64
	}

	sortedFaces := make([]sortableFace, len(fs.faces))
	for i, f := range fs.faces {
		sortedFaces[i] = sortableFace{
			f:      f,
			distSq: f.GetMidPoint().DistanceSquaredTo(pos),
		}
	}

	// Sort the faces by distance to the position
	sort.Slice(sortedFaces, func(i, j int) bool {
		return sortedFaces[i].distSq > sortedFaces[j].distSq
	})

	for i, sf := range sortedFaces {
		fs.faces[i] = sf.f
	}
}
