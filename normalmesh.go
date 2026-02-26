package gosie3d

func NewNormalMesh() *NormalMesh {
	return &NormalMesh{Mesh: *NewMesh()}
}

type NormalMesh struct {
	Mesh
}

func (nm *NormalMesh) AddNormal(pnts Vector3) (Vector3, int) {
	// AddPoint takes a Vector3.
	// Since NormalMesh embeds Mesh, we can call AddPoint directly.
	return nm.AddPoint(pnts)
}

func (nm *NormalMesh) Copy() *NormalMesh {
	return &NormalMesh{Mesh: *nm.Mesh.Copy()}
}
