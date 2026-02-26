package universe

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/smasonuk/gosie3d"
)

func init() {
	GalaxyStars = GenerateSpiralGalaxy(300000, 1772054134190328000)
}

type Starfield struct {
	GalaxyStars *Galaxy
	Camera      *gosie3d.Camera
	Position    gosie3d.Vector3
}

func NewStarfield(camera *gosie3d.Camera, pos gosie3d.Vector3) *Starfield {
	return &Starfield{
		GalaxyStars: GalaxyStars,
		Camera:      camera,
		Position:    pos,
	}
}

type GalacticStar struct {
	Position   gosie3d.Vector3
	Luminosity float64
	BaseColor  color.RGBA
	IsGas      bool
	IsDust     bool
}

type Galaxy struct {
	Stars []GalacticStar
}

var GalaxyStars *Galaxy

// creates a procedural galaxy with a core and spiral arms.
func GenerateSpiralGalaxy(totalStars int, seed int64) *Galaxy {
	rand.New(rand.NewSource(seed))

	galaxy := &Galaxy{
		Stars: make([]GalacticStar, 0, totalStars),
	}

	// Galaxy Parameters (Tweak these to change the galaxy's shape!)
	const (
		numArms       = 2       // Most classic spirals have 2 major arms
		armWrap       = 5.0     // How many radians the arms twist from core to edge
		maxRadius     = 50000.0 // Radius of the galaxy in Light-Years
		coreRadius    = 8000.0  // Radius of the dense central bulge
		diskThickness = 2000.0  // How "thick" the flat disk is
		coreRatio     = 0.3     // 30% of stars go in the core, 70% in the arms
	)

	numCoreStars := int(float64(totalStars) * coreRatio)
	numArmStars := totalStars - numCoreStars

	// 1. generate the galactic core
	for i := 0; i < numCoreStars; i++ {
		// Use normal (Gaussian) distribution to cluster stars tightly in the center
		x := rand.NormFloat64() * (coreRadius / 3.0)
		y := rand.NormFloat64() * (coreRadius / 3.0)
		z := rand.NormFloat64() * (coreRadius / 3.0)

		// Core stars are generally older, yellower, and have lower average luminosity
		// lum := math.Pow(rand.Float64(), 4.0) * 100.0 // Heavy bias towards dim stars
		lum := math.Pow(rand.Float64(), 4.0) * 2000.0

		galaxy.Stars = append(galaxy.Stars, GalacticStar{
			Position:   gosie3d.NewVector3(x, y, z),
			Luminosity: lum + 5.0,                      // Base minimum luminosity
			BaseColor:  color.RGBA{255, 230, 200, 255}, // Warm yellow-white
		})
	}

	// 2. generate the spiral arms
	for i := 0; i < numArmStars; i++ {
		// Pick a random distance from the center.
		// We use a square root to ensure even distribution across the disk area.
		dist := math.Sqrt(rand.Float64()) * maxRadius

		// Which arm does this star belong to?
		armIndex := rand.Intn(numArms)
		armOffset := (float64(armIndex) / float64(numArms)) * 2.0 * math.Pi

		// Calculate the base spiral angle
		spiralAngle := armWrap * (dist / maxRadius)

		// Add organic "fuzziness" to the arms.
		// The spread gets wider the further out from the core you go.
		spreadSpread := (dist / maxRadius) * 0.5
		randomAngleOffset := rand.NormFloat64() * spreadSpread

		// Final angle calculation
		theta := spiralAngle + armOffset + randomAngleOffset

		// Convert polar coordinates to Cartesian (X, Z)
		x := dist * math.Cos(theta)
		z := dist * math.Sin(theta)

		// Calculate Y (height). The disk gets slightly thicker at the edges.
		thicknessAtDist := diskThickness * (1.0 + (dist / maxRadius))
		y := rand.NormFloat64() * (thicknessAtDist / 4.0)

		// Arm stars are younger, bluer, and feature rare, incredibly bright super-giants
		// lum := math.Pow(rand.Float64(), 6.0) * 10000.0 // Bias towards dim, but allows massive spikes
		// Arm stars are younger, bluer, and feature rare, incredibly bright super-giants
		lum := math.Pow(rand.Float64(), 6.0) * 10000.0

		isGas := rand.Float64() < 0.15 // 15% chance to be a gas cloud instead of a star

		var starColor color.RGBA
		var starLum float64

		if isGas {
			// Nebulae glow in bright pinks, purples, and cyans (H-alpha and Oxygen emissions)
			if rand.Float64() > 0.5 {
				starColor = color.RGBA{220, 50, 150, 255} // Pink/Magenta
			} else {
				starColor = color.RGBA{50, 200, 250, 255} // Cyan
			}
			starLum = lum * 1.5 // Gas clouds are bright but diffuse
		} else {
			starColor = color.RGBA{200, 220, 255, 255} // Cool blue-white
			starLum = lum + 10.0
		}

		galaxy.Stars = append(galaxy.Stars, GalacticStar{
			Position:   gosie3d.NewVector3(x, y, z),
			Luminosity: starLum,
			BaseColor:  starColor,
			IsGas:      isGas, // Flag it!
		})
	}

	// generate dark dust lanes
	const dustRatio = 0.5 // We need a LOT of dust to block the light
	numDustClouds := int(float64(totalStars) * dustRatio)

	for i := 0; i < numDustClouds; i++ {
		dist := math.Sqrt(rand.Float64()) * maxRadius

		armIndex := rand.Intn(numArms)
		armOffset := (float64(armIndex) / float64(numArms)) * 2.0 * math.Pi
		spiralAngle := armWrap * (dist / maxRadius)

		// dust placement: Offset the angle slightly backwards so it hugs the inside of the arms
		insideEdgeOffset := -0.15
		spread := (dist / maxRadius) * 0.2 // Tighter spread than the stars

		theta := spiralAngle + armOffset + insideEdgeOffset + (rand.NormFloat64() * spread)

		x := dist * math.Cos(theta)
		z := dist * math.Sin(theta)

		// Dust is extremely flat compared to the rest of the disk
		y := rand.NormFloat64() * (diskThickness / 8.0)

		galaxy.Stars = append(galaxy.Stars, GalacticStar{
			Position:   gosie3d.NewVector3(x, y, z),
			Luminosity: rand.Float64() * 8000.0, // Acts as "darkness" strength
			IsDust:     true,
		})
	}

	// temporarily cranked up to 15% so you can see it!
	const haloRatio = 0.15
	numHaloStars := int(float64(totalStars) * haloRatio)

	// GENERATE THE GALACTIC HALO
	for i := 0; i < numHaloStars; i++ {
		x := rand.NormFloat64() * maxRadius
		y := rand.NormFloat64() * maxRadius
		z := rand.NormFloat64() * maxRadius

		// Crank the luminosity up temporarily so they burn brightly into the sensor
		lum := rand.Float64() * 5000.0 * 4

		galaxy.Stars = append(galaxy.Stars, GalacticStar{
			Position:   gosie3d.NewVector3(x, y, z),
			Luminosity: lum + 100.0,
			BaseColor:  color.RGBA{255, 150, 150, 255}, // Make them distinctly red to stand out
		})
	}

	return galaxy
}

func (g *Galaxy) TakeProbeSnapshot(cam *gosie3d.Camera, probeGalacticPos gosie3d.Vector3, width, height int, exposure float64) *image.RGBA {

	// 1. 3-Channel Sensor Array (Red, Green, Blue)
	sensorR := make([][]float64, width)
	sensorG := make([][]float64, width)
	sensorB := make([][]float64, width)
	for i := range sensorR {
		sensorR[i] = make([]float64, height)
		sensorG[i] = make([]float64, height)
		sensorB[i] = make([]float64, height)
	}

	viewMat := cam.GetMatrix()

	// 2. Accumulate light (Photons)
	for _, star := range g.Stars {
		relPos := gosie3d.Subtract(star.Position, probeGalacticPos)
		distSq := relPos.X*relPos.X + relPos.Y*relPos.Y + relPos.Z*relPos.Z
		if distSq < 1.0 {
			distSq = 1.0
		}

		apparentBrightness := (star.Luminosity / distSq) * exposure

		dist := math.Sqrt(distSq)
		dir := gosie3d.NewVector3(relPos.X/dist, relPos.Y/dist, relPos.Z/dist)
		camSpaceDir := viewMat.RotateVector3(dir)

		if camSpaceDir.Z <= 0 {
			continue
		}

		screenX := int(gosie3d.ConvertToScreenX(float64(width), float64(height), camSpaceDir.X, camSpaceDir.Z))
		screenY := int(gosie3d.ConvertToScreenY(float64(width), float64(height), camSpaceDir.Y, camSpaceDir.Z))

		// Ensure we have padding for splatting and flares
		if screenX >= 15 && screenX < width-15 && screenY >= 15 && screenY < height-15 {

			r := float64(star.BaseColor.R) / 255.0
			g_col := float64(star.BaseColor.G) / 255.0
			b := float64(star.BaseColor.B) / 255.0

			if star.IsDust {
				// NEGATIVE LIGHT (Dust Lane)
				for dx := -2; dx <= 2; dx++ {
					for dy := -2; dy <= 2; dy++ {
						distSqSplat := float64(dx*dx + dy*dy)
						if distSqSplat > 4.0 {
							continue
						}

						weight := 1.0 / (distSqSplat + 1.0)
						darkness := (apparentBrightness * 0.4) * weight

						sensorR[screenX+dx][screenY+dy] -= darkness
						sensorG[screenX+dx][screenY+dy] -= darkness
						sensorB[screenX+dx][screenY+dy] -= darkness
					}
				}
			} else if star.IsGas {
				// WIDE SPLAT (Nebula Cloud)
				for dx := -3; dx <= 3; dx++ {
					for dy := -3; dy <= 3; dy++ {
						distSqSplat := float64(dx*dx + dy*dy)
						if distSqSplat > 9.0 {
							continue
						}

						weight := 1.0 / (distSqSplat + 1.0)
						diffuseLight := (apparentBrightness * 0.15) * weight

						sensorR[screenX+dx][screenY+dy] += diffuseLight * r
						sensorG[screenX+dx][screenY+dy] += diffuseLight * g_col
						sensorB[screenX+dx][screenY+dy] += diffuseLight * b
					}
				}
			} else {
				// TIGHT SPLAT (Normal Star)
				center := apparentBrightness * 0.6
				side := apparentBrightness * 0.1

				sensorR[screenX][screenY] += center * r
				sensorG[screenX][screenY] += center * g_col
				sensorB[screenX][screenY] += center * b

				offsets := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
				for _, off := range offsets {
					sensorR[screenX+off[0]][screenY+off[1]] += side * r
					sensorG[screenX+off[0]][screenY+off[1]] += side * g_col
					sensorB[screenX+off[0]][screenY+off[1]] += side * b
				}

				// If the star is incredibly bright, it creates a cross flare on the lens
				if apparentBrightness > 15.0 {
					// The brighter the star, the longer the spike (capped at 12 pixels)
					spikeLen := int(math.Min(12.0, apparentBrightness/3.0))
					spikeStrength := apparentBrightness * 0.05

					for s := 1; s <= spikeLen; s++ {
						// Fade the light out towards the tips of the spike
						fade := 1.0 - (float64(s) / float64(spikeLen))
						lightVal := spikeStrength * fade

						spikeOffsets := [][]int{{s, 0}, {-s, 0}, {0, s}, {0, -s}}
						for _, off := range spikeOffsets {
							sensorR[screenX+off[0]][screenY+off[1]] += lightVal * r
							sensorG[screenX+off[0]][screenY+off[1]] += lightVal * g_col
							sensorB[screenX+off[0]][screenY+off[1]] += lightVal * b
						}
					}
				}
			}
		}
	}

	// 3. Develop with Noise, Tone Mapping, and Scanlines
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Add "Sensor Noise" (grain)
			noise := (rand.Float64() - 0.5) * 0.008

			sR := math.Max(0, sensorR[x][y])
			sG := math.Max(0, sensorG[x][y])
			sB := math.Max(0, sensorB[x][y])

			finalR := sR / (1.0 + sR)
			finalG := sG / (1.0 + sG)
			finalB := sB / (1.0 + sB)

			// Space background tint
			bgR, bgG, bgB := 0.02, 0.02, 0.03

			r := clamp(int((finalR+bgR+noise)*255), 0, 255)
			g_c := clamp(int((finalG+bgG+noise)*255), 0, 255)
			b := clamp(int((finalB+bgB+noise)*255), 0, 255)

			// Darken every even horizontal row by 20% to simulate a telemetry feed
			if y%2 == 0 {
				r = int(float64(r) * 0.8)
				g_c = int(float64(g_c) * 0.8)
				b = int(float64(b) * 0.8)
			}

			img.Set(x, y, color.RGBA{uint8(r), uint8(g_c), uint8(b), 255})
		}
	}

	return img
}

// Helper function
func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (s *Starfield) GetStarField(height, width int) *image.RGBA {
	galaxy := GalaxyStars

	// Exposure of 50k - 100k is usually the "sweet spot" for this distance
	snapshot := galaxy.TakeProbeSnapshot(
		s.Camera,
		s.Position,
		width,
		height,
		80000.0)
	return snapshot
}
