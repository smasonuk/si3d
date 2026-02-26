package si3d

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	imgdraw "image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "regenerate golden test images")

const (
	goldenWidth    = 200
	goldenHeight   = 150
	goldenDir      = "../../testdata"
	goldenFilename = "golden_render.png"

	goldenMountainsWidth    = 400
	goldenMountainsHeight   = 300
	goldenMountainsFilename = "golden_render_mountains.png"
)

func buildGoldenScene() *image.RGBA {
	model := NewRectangle(80, 60, 80, color.RGBA{R: 100, G: 150, B: 200, A: 255})
	world := NewWorld3d()

	cam := NewCamera(0, 0, 0, 0, 0, 0)
	cam.SetCameraPosition(200, -120, -350)
	cam.LookAt(NewVector3(0, 0, 0), NewVector3(0, -1, 0))
	world.AddCamera(cam, 200, -120, -350)
	world.AddObject(&Entity{Model: model, X: 0, Y: 0, Z: 0})

	bgColor := color.RGBA{R: 20, G: 20, B: 40, A: 255}
	return world.Render(goldenWidth, goldenHeight, bgColor)
}

func loadGoldenImage(t *testing.T, path string) *image.RGBA {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open golden image %s: %v", path, err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode golden image: %v", err)
	}
	rgba := image.NewRGBA(img.Bounds())
	imgdraw.Draw(rgba, rgba.Bounds(), img, image.Point{}, imgdraw.Src)
	return rgba
}

func saveGoldenImage(t *testing.T, path string, img *image.RGBA) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create testdata dir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create golden image: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode golden image: %v", err)
	}
	t.Logf("wrote golden image to %s", path)
}

func compareImages(t *testing.T, got, want *image.RGBA) bool {
	t.Helper()
	if got.Bounds() != want.Bounds() {
		t.Errorf("image size mismatch: got %v, want %v", got.Bounds(), want.Bounds())
		return false
	}
	if bytes.Equal(got.Pix, want.Pix) {
		return true
	}
	w := got.Bounds().Dx()
	for i := 0; i < len(got.Pix); i += 4 {
		if got.Pix[i] != want.Pix[i] || got.Pix[i+1] != want.Pix[i+1] ||
			got.Pix[i+2] != want.Pix[i+2] || got.Pix[i+3] != want.Pix[i+3] {
			idx := i / 4
			t.Errorf("pixel mismatch at (%d,%d): got RGBA(%d,%d,%d,%d) want RGBA(%d,%d,%d,%d)",
				idx%w, idx/w,
				got.Pix[i], got.Pix[i+1], got.Pix[i+2], got.Pix[i+3],
				want.Pix[i], want.Pix[i+1], want.Pix[i+2], want.Pix[i+3])
			break
		}
	}
	return false
}

// TestRenderGolden is the main E2E golden image regression test.
// Normal run: go test -run TestRenderGolden
// Regenerate:  go test -run TestRenderGolden -update
func TestRenderGolden(t *testing.T) {
	goldenPath := filepath.Join(goldenDir, goldenFilename)
	got := buildGoldenScene()

	if *update {
		saveGoldenImage(t, goldenPath, got)
		return
	}

	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		t.Fatalf("golden image not found at %s\nRun: go test -run TestRenderGolden -update", goldenPath)
	}

	want := loadGoldenImage(t, goldenPath)
	if !compareImages(t, got, want) {
		failPath := filepath.Join(goldenDir, "golden_render_FAIL.png")
		saveGoldenImage(t, failPath, got)
		t.Logf("failing render saved to %s", failPath)
		t.Fail()
	}
}

// buildMountainsScene constructs a deterministic mountains scene using Perlin noise
// with a fixed seed, mirroring cmd/main.go but at a smaller size for test speed.
func buildMountainsScene() *image.RGBA {
	mountains := NewSubdividedPlaneHeightMapPerlin(
		10000, 10000,
		color.RGBA{R: 153, G: 196, B: 210, A: 255},
		15, 800, 800, 42,
	)
	mountains.SetDrawLinesOnly(true)

	cam := NewCamera(0, 0, 0, 0, 0, 0)
	camX := 500.0
	camZ := 0.0
	cam.SetCameraPosition(camX, -200, camZ)
	cam.LookAt(NewVector3(0, -100, 0), NewVector3(0, -1, 0))

	world := NewWorld3d()
	world.AddCamera(cam, camX, -200, camZ)
	world.AddObjectDrawFirst(&Entity{Model: mountains, X: 0, Y: 0, Z: 0})

	bgColor := color.RGBA{R: 10, G: 10, B: 30, A: 255}
	return world.Render(goldenMountainsWidth, goldenMountainsHeight, bgColor)
}

// TestRenderGoldenMountains is an E2E golden image regression test for the mountains scene.
// Normal run: go test -run TestRenderGoldenMountains
// Regenerate:  go test -run TestRenderGoldenMountains -update
func TestRenderGoldenMountains(t *testing.T) {
	goldenPath := filepath.Join(goldenDir, goldenMountainsFilename)
	got := buildMountainsScene()

	if *update {
		saveGoldenImage(t, goldenPath, got)
		return
	}

	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		t.Fatalf("golden image not found at %s\nRun: go test -run TestRenderGoldenMountains -update", goldenPath)
	}

	want := loadGoldenImage(t, goldenPath)
	if !compareImages(t, got, want) {
		failPath := filepath.Join(goldenDir, "golden_render_mountains_FAIL.png")
		saveGoldenImage(t, failPath, got)
		t.Logf("failing render saved to %s", failPath)
		t.Fail()
	}
}

// TestRenderGoldenMountainsSanity checks that the mountains render is non-blank.
// This test does not require the golden image file.
func TestRenderGoldenMountainsSanity(t *testing.T) {
	got := buildMountainsScene()
	bgColor := color.RGBA{R: 10, G: 10, B: 30, A: 255}
	bgR, bgG, bgB, bgA := bgColor.R, bgColor.G, bgColor.B, bgColor.A
	for i := 0; i < len(got.Pix); i += 4 {
		if got.Pix[i] != bgR || got.Pix[i+1] != bgG ||
			got.Pix[i+2] != bgB || got.Pix[i+3] != bgA {
			return
		}
	}
	t.Error("mountains render produced an entirely blank image - no geometry was drawn")
}

// TestRenderGoldenSanity checks that the render is non-blank (geometry is visible).
// This test does not require the golden image file.
func TestRenderGoldenSanity(t *testing.T) {
	got := buildGoldenScene()
	bgColor := color.RGBA{R: 20, G: 20, B: 40, A: 255}
	bgR, bgG, bgB, bgA := bgColor.R, bgColor.G, bgColor.B, bgColor.A
	for i := 0; i < len(got.Pix); i += 4 {
		if got.Pix[i] != bgR || got.Pix[i+1] != bgG ||
			got.Pix[i+2] != bgB || got.Pix[i+3] != bgA {
			return // found a non-background pixel, test passes
		}
	}
	t.Error("render produced an entirely blank image - no geometry was drawn")
}
