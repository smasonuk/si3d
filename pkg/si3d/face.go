package si3d

import "image/color"

type Face struct {
	Points    []Vector3
	Col       color.RGBA
	normal    Vector3
	hasNormal bool
	plane     *Plane
	vecPnts   []Vector3
	meRev     bool
	Cnum      int
}

const (
	FACE_NORMAL  = 0
	FACE_REVERSE = 1
)

func NewFaceEmpty(col color.RGBA, normal Vector3) *Face {
	has := normal != Vector3{}
	return &Face{
		Points:    make([]Vector3, 0),
		Col:       col,
		normal:    normal,
		hasNormal: has,
	}
}

// copy
func (f *Face) Copy() *Face {
	newFace := &Face{
		Points:    make([]Vector3, len(f.Points)),
		Col:       f.Col,
		normal:    f.GetNormal(),
		hasNormal: true,
		plane:     f.plane,
		vecPnts:   make([]Vector3, len(f.vecPnts)),
		meRev:     f.meRev,
		Cnum:      f.Cnum,
	}
	copy(newFace.Points, f.Points)
	copy(newFace.vecPnts, f.vecPnts)
	return newFace
}

func NewFace(pnts []Vector3, col color.RGBA, normal Vector3) *Face {
	has := normal != Vector3{}
	f := &Face{
		Points:    pnts,
		Col:       col,
		normal:    normal,
		hasNormal: has,
	}
	if pnts != nil {
		f.Cnum = len(pnts)
	}
	return f
}

func (f *Face) SetColor(col color.RGBA) {
	f.Col = col
}

func (f *Face) AddPoint(x, y, z float64) {
	pnts := NewVector3(x, y, z)
	f.vecPnts = append(f.vecPnts, pnts)
}

func (f *Face) GetNormal() Vector3 {
	if !f.hasNormal {
		f.createNormal()
	}
	return f.normal
}

func (f *Face) SetNormal(norm Vector3) {
	f.normal = norm
	f.hasNormal = true
}

func (f *Face) Finished(reverse int) {
	if reverse == FACE_REVERSE {
		f.meRev = true
	}
	f.Cnum = len(f.vecPnts)
	f.Points = make([]Vector3, f.Cnum)
	copy(f.Points, f.vecPnts)
	f.vecPnts = nil
}

func (f *Face) GetPlane() *Plane {
	if f.plane != nil {
		return f.plane
	}
	f.plane = NewPlane(f, f.GetNormal())
	return f.plane
}

func (f *Face) createNormal() {
	if len(f.Points) < 3 {
		f.normal = NewVector3(0, 0, 1)
		f.hasNormal = true
		return
	}

	x1, y1, z1 := f.Points[0].X, f.Points[0].Y, f.Points[0].Z
	x2, y2, z2 := f.Points[1].X, f.Points[1].Y, f.Points[1].Z
	x3, y3, z3 := f.Points[2].X, f.Points[2].Y, f.Points[2].Z

	u1, u2, u3 := x2-x1, y2-y1, z2-z1
	v1, v2, v3 := x3-x2, y3-y2, z3-z2

	f.normal = Vector3{}
	if f.meRev {
		f.normal.X = -(u2*v3 - u3*v2)
		f.normal.Y = -(u3*v1 - u1*v3)
		f.normal.Z = -(u1*v2 - u2*v1)
	} else {
		f.normal.X = u2*v3 - u3*v2
		f.normal.Y = u3*v1 - u1*v3
		f.normal.Z = u1*v2 - u2*v1
	}

	// Normalizing
	f.normal = f.normal.Normalize()
	f.hasNormal = true
}

// get midpoint of the face
func (f *Face) GetMidPoint() Vector3 {
	if len(f.Points) == 0 {
		return NewVector3(0, 0, 0)
	}

	sumX, sumY, sumZ := 0.0, 0.0, 0.0
	for _, p := range f.Points {
		sumX += p.X
		sumY += p.Y
		sumZ += p.Z
	}
	count := float64(len(f.Points))
	return NewVector3(sumX/count, sumY/count, sumZ/count)
}

// get the distance from the face to a point
func (f *Face) GetDistanceToPoint(p Vector3) float64 {
	midpoint := f.GetMidPoint()
	return midpoint.DistanceTo(p)
}
