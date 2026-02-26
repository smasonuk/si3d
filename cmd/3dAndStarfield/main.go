package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"os"

	"github.com/smasonuk/gosie3d"
	"github.com/smasonuk/gosie3d/universe"
)

type Probe struct {
	StarfieldPosition gosie3d.Vector3 // Probe's location in the galaxy (Light-Years)
	Camera            *gosie3d.Camera // The Master Camera
}

func NewProbe(starfieldPosition gosie3d.Vector3, cam *gosie3d.Camera) *Probe {
	return &Probe{
		StarfieldPosition: starfieldPosition,
		Camera:            cam,
	}
}

// func main() {
// 	frames := 60
// 	outGif := &gif.GIF{}

// 	// 1. Where is the probe in the GALAXY?
// 	starfieldPosition := gosie3d.NewVector3(10000, 25000, 35000)

// 	// 2. Pre-generate the Mountains so we don't rebuild the heightmap 60 times
// 	fmt.Println("Generating mountains...")
// 	mountains := gosie3d.NewSubdividedPlaneHeightMapPerlin(
// 		10000, 10000,
// 		color.RGBA{R: 153, G: 196, B: 210, A: 255},
// 		35, 800, 800, 42,
// 	)
// 	mountains.SetDrawLinesOnly(true)

// 	fmt.Println("Rendering animation frames...")

// 	for i := 0; i < frames; i++ {
// 		// 3. Calculate a sweeping 360-degree orbit around the mountains
// 		progress := float64(i) / float64(frames)
// 		// cameraAngle := progress * (math.Pi * 2.0)

// 		// cameraDistance := 500.0
// 		// cameraHeight := -200.0

// 		// localCamX := math.Cos(cameraAngle) * cameraDistance
// 		// localCamZ := math.Sin(cameraAngle) * cameraDistance

// 		// // 4. Create the Master Camera for this specific frame
// 		// masterCam := gosie3d.NewCamera(localCamX, cameraHeight, localCamZ, 0, 0, 0)
// 		// masterCam.LookAt(gosie3d.NewVector3(0, -100, 0), gosie3d.NewVector3(0, 1, 0))

// 		// 1. Stand in one spot (Fixed Position)
// 		localCamX := 0.0
// 		cameraHeight := -200.0
// 		localCamZ := -800.0 // Backed away from the mountains

// 		// 2. Turn our head (Pan left to right)
// 		// Sweeps from -45 degrees to +45 degrees over the 60 frames
// 		panAngle := (progress - 0.5) * (math.Pi / 2.0)

// 		// 3. Calculate a target to look at that sweeps across the horizon
// 		lookAtX := localCamX + (math.Sin(panAngle) * 100.0)
// 		lookAtZ := localCamZ + (math.Cos(panAngle) * 100.0)

// 		// 4. Create the Master Camera for this specific frame
// 		masterCam := gosie3d.NewCamera(localCamX, cameraHeight, localCamZ, 0, 0, 0)

// 		// Look slightly down at the mountains while panning
// 		masterCam.LookAt(gosie3d.NewVector3(lookAtX, cameraHeight-20.0, lookAtZ), gosie3d.NewVector3(0, 1, 0))

// 		probe := NewProbe(starfieldPosition, masterCam)

// 		// 5. Render the Starfield Background
// 		field := universe.NewStarfield(probe.Camera, probe.StarfieldPosition)
// 		frameImg := field.GetStarField(512, 512) // Specify the width and height for the starfield image

// 		// 6. Render the Mountains on top
// 		world := gosie3d.NewWorld3d()
// 		world.AddCamera(probe.Camera, localCamX, cameraHeight, localCamZ)
// 		world.AddObjectDrawFirst(&gosie3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})
// 		world.RenderToImage(frameImg)

// 		// 7. Convert the RGBA frame to a 256-color Paletted image for the GIF
// 		bounds := frameImg.Bounds()
// 		palettedImage := image.NewPaletted(bounds, palette.Plan9)
// 		draw.Draw(palettedImage, palettedImage.Rect, frameImg, bounds.Min, draw.Over)

// 		// 8. Append to the GIF sequence
// 		outGif.Image = append(outGif.Image, palettedImage)
// 		outGif.Delay = append(outGif.Delay, 5) // 5 = 50ms delay per frame (20fps)

// 		fmt.Printf("Rendered frame %d/%d\n", i+1, frames)
// 	}

// 	// 9. Save the final animated GIF
// 	fmt.Println("Saving orbit.gif...")
// 	f, err := os.Create("orbit.gif")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	gif.EncodeAll(f, outGif)
// 	fmt.Println("Done! Open orbit.gif to see the result.")
// }

func main() {
	// Bumped to 120 frames for a smoother, slightly slower 360-degree flight
	frames := 120
	outGif := &gif.GIF{}

	// 1. Where is the probe in the GALAXY?
	starfieldPosition := gosie3d.NewVector3(10000, 25000, 35000)

	// 2. Pre-generate the Mountains so we don't rebuild the heightmap every frame
	fmt.Println("Generating mountains...")
	mountains := gosie3d.NewSubdividedPlaneHeightMapPerlin(
		10000, 10000,
		color.RGBA{R: 153, G: 196, B: 210, A: 255},
		35, 800, 800, 42,
	)
	mountains.SetDrawLinesOnly(true)

	fmt.Println("Rendering animation frames...")

	for i := 0; i < frames; i++ {
		// 3. Calculate progression (0.0 to 1.0) and convert to radians
		progress := float64(i) / float64(frames)
		theta := progress * (math.Pi * 2.0) // Full 360 circle

		circleRadius := 400.0 // The size of our "smallish" circular path
		cameraHeight := -200.0

		// 4. Calculate camera POSITION on the circle
		localCamX := math.Cos(theta) * circleRadius
		localCamZ := math.Sin(theta) * circleRadius

		// 5. Calculate the TANGENT vector (the direction straight ahead on the curve)
		// The derivative of (cos, sin) is (-sin, cos)
		tangentX := -math.Sin(theta)
		tangentZ := math.Cos(theta)

		// 6. Project the LookAt target out along the tangent vector
		lookDistance := 100.0
		lookAtX := localCamX + (tangentX * lookDistance)
		lookAtZ := localCamZ + (tangentZ * lookDistance)
		lookAtY := cameraHeight - 20.0 // Look slightly down at the passing mountains

		// 7. Create the Master Camera for this specific frame
		masterCam := gosie3d.NewCamera(localCamX, cameraHeight, localCamZ, 0, 0, 0)
		masterCam.LookAt(gosie3d.NewVector3(lookAtX, lookAtY, lookAtZ), gosie3d.NewVector3(0, 1, 0))

		probe := NewProbe(starfieldPosition, masterCam)

		// 8. Render the Starfield Background
		field := universe.NewStarfield(probe.Camera, probe.StarfieldPosition)
		frameImg := field.GetStarField(512, 512) // Assuming GetStarField doesn't need w/h parameters if hardcoded

		// 9. Render the Mountains on top
		world := gosie3d.NewWorld3d()
		world.AddCamera(probe.Camera, localCamX, cameraHeight, localCamZ)
		world.AddObjectDrawFirst(&gosie3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})
		world.RenderToImage(frameImg)

		// 10. Convert the RGBA frame to a 256-color Paletted image for the GIF
		bounds := frameImg.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, frameImg, bounds.Min, draw.Over)

		// 11. Append to the GIF sequence
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 5) // 5 = 50ms delay per frame (20fps)

		fmt.Printf("Rendered frame %d/%d\n", i+1, frames)
	}

	// 12. Save the final animated GIF
	fmt.Println("Saving orbit.gif...")
	f, err := os.Create("orbit.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gif.EncodeAll(f, outGif)
	fmt.Println("Done! Open orbit.gif to see the result.")
}
