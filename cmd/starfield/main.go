package main

import (
	"image/png"
	"os"

	"github.com/smasonuk/gosie3d"
	"github.com/smasonuk/gosie3d/universe"
)

func main() {
	probePos := gosie3d.NewVector3(10000, 25000, 35000)
	cam := gosie3d.NewCamera(probePos.X, probePos.Y, probePos.Z, 0, 0, 0)
	cam.LookAt(gosie3d.NewVector3(0, 0, 0), gosie3d.NewVector3(0, 1, 0))

	field := universe.NewStarfield(cam, probePos)
	snapshot := field.GetStarField(512, 512)

	filename := ".temp.png"
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	png.Encode(f, snapshot)
	f.Close()

}
