


# Examples
## Create a animation
```go
func main() {
	// I bumped the star count to 300,000 for a slightly denser look during motion
	// galaxy := GenerateSpiralGalaxy(300000, 1772054134190328000)
	frames := 50
	fmt.Printf("Starting render of %d frames...\n", frames)
	for i := 0; i < frames; i++ {
		// Calculate the probe's path (A sweeping orbit)
		// We sweep through 90 degrees (Pi/2) over the 60 frames
		progress := float64(i) / float64(frames)
		angle := progress * (math.Pi / 2.0)
		orbitRadius := 40000.0
		// Probe moves in a wide arc on the X/Z plane
		camX := math.Cos(angle) * orbitRadius
		camZ := math.Sin(angle) * orbitRadius
		// The probe slowly descends towards the galactic plane
		camY := 25000.0 - (progress * 15000.0)
		probePos := gosie3d.NewVector3(camX, camY, camZ)
		cam := gosie3d.NewCamera(probePos.X, probePos.Y, probePos.Z, 0, 0, 0)
		// Always keep the camera locked onto the galactic core
		cam.LookAt(gosie3d.NewVector3(0, 0, 0), gosie3d.NewVector3(0, 1, 0))
		// Render the frame!
		snapshot := galaxy.TakeProbeSnapshot(cam, probePos, 512, 512, 80000.0)
		// Save sequentially numbered files
		filename := fmt.Sprintf("./images/frame_%03d.png", i)
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		png.Encode(f, snapshot)
		f.Close()
		// Print a progress indicator so you know the CPU is working
		fmt.Printf("Rendered %s (%.1f%% complete)\n", filename, (progress * 100.0))
	}
	fmt.Println("Render complete! You can now stitch these into a GIF/Video.")
}
```




```

// func main() {

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

// 	// 3. Calculate a sweeping 360-degree orbit around the mountains
// 	progress := 0.0
// 	// cameraAngle := progress * (math.Pi * 2.0)

// 	// cameraDistance := 500.0
// 	// cameraHeight := -200.0

// 	// localCamX := math.Cos(cameraAngle) * cameraDistance
// 	// localCamZ := math.Sin(cameraAngle) * cameraDistance

// 	// // 4. Create the Master Camera for this specific frame
// 	// masterCam := gosie3d.NewCamera(localCamX, cameraHeight, localCamZ, 0, 0, 0)
// 	// masterCam.LookAt(gosie3d.NewVector3(0, -100, 0), gosie3d.NewVector3(0, 1, 0))

// 	// 1. Stand in one spot (Fixed Position)
// 	localCamX := 0.0
// 	cameraHeight := -200.0
// 	localCamZ := -800.0 // Backed away from the mountains

// 	// 2. Turn our head (Pan left to right)
// 	// Sweeps from -45 degrees to +45 degrees over the 60 frames
// 	panAngle := (progress - 0.5) * (math.Pi / 2.0)

// 	// 3. Calculate a target to look at that sweeps across the horizon
// 	lookAtX := localCamX + (math.Sin(panAngle) * 100.0)
// 	lookAtZ := localCamZ + (math.Cos(panAngle) * 100.0)

// 	// 4. Create the Master Camera for this specific frame
// 	masterCam := gosie3d.NewCamera(localCamX, cameraHeight, localCamZ, 0, 0, 0)

// 	// Look slightly down at the mountains while panning
// 	masterCam.LookAt(gosie3d.NewVector3(lookAtX, cameraHeight-20.0, lookAtZ), gosie3d.NewVector3(0, 1, 0))

// 	probe := NewProbe(starfieldPosition, masterCam)

// 	// 5. Render the Starfield Background
// 	field := universe.NewStarfield(probe.Camera, probe.StarfieldPosition)
// 	frameImg := field.GetStarField(512, 512) // Specify the width and height for the starfield image

// 	// 6. Render the Mountains on top
// 	world := gosie3d.NewWorld3d()
// 	world.AddCamera(probe.Camera, localCamX, cameraHeight, localCamZ)
// 	world.AddObjectDrawFirst(&gosie3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})
// 	world.RenderToImage(frameImg)

// 	// 7. Convert the RGBA frame to a 256-color Paletted image for the GIF
// 	bounds := frameImg.Bounds()
// 	palettedImage := image.NewPaletted(bounds, palette.Plan9)
// 	draw.Draw(palettedImage, palettedImage.Rect, frameImg, bounds.Min, draw.Over)

// 	// save to image
// 	f, err := os.Create("frame.png")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	gif.Encode(f, palettedImage, nil)

// 	fmt.Println("Done! Open frame.png to see the result.")

// }
```