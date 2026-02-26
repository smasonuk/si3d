package main

import (
	"image/color"
	"log"
	"math"

	"github.com/smasonuk/si3d/pkg/si3d"
)

const (
	WIDTH  = 512
	HEIGHT = 512
)

func main() {
	// Create the mountains heightmap
	mountains := si3d.NewSubdividedPlaneHeightMapPerlin(
		10000, 10000,
		color.RGBA{R: 153, G: 196, B: 210, A: 255},
		35, 800, 800, 42,
	)
	mountains.SetDrawLinesOnly(true)

	// Create camera
	cam := si3d.NewCamera(0, 0, 0, 0, 0, 0)

	cameraAngle := 0.0
	cameraHeight := -200.0
	cameraDistance := 500.0

	camX := math.Cos(cameraAngle) * cameraDistance
	camZ := math.Sin(cameraAngle) * cameraDistance

	cam.SetCameraPosition(camX, cameraHeight, camZ)
	cam.LookAt(
		si3d.NewVector3(0, -100, 0),
		si3d.NewVector3(0, -1, 0),
	)

	// Build the world
	world := si3d.NewWorld3d()
	world.AddCamera(cam, cam.GetPosition().X, cam.GetPosition().Y, cam.GetPosition().Z)
	world.AddObjectDrawFirst(&si3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})

	// Render to file
	bgColor := color.RGBA{R: 10, G: 10, B: 30, A: 255}
	err := world.RenderToFile(WIDTH, HEIGHT, bgColor, "output.png")
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Rendered to output.png")
}
