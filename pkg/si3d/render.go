package si3d

import (
	"image"
	"image/color"
	imgdraw "image/draw"
	"image/png"
	"os"
)

// Render draws the current state of the world into a new *image.RGBA
// of the given dimensions.
// The background is filled with bgColor before any 3D objects are drawn.
func (w *World) Render(width, height int, bgColor color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fast background fill
	imgdraw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.ZP, imgdraw.Src)

	w.PaintObjects(img, width, height)
	return img
}

func (w *World) RenderToImage(img *image.RGBA) *image.RGBA {
	// img := image.NewRGBA(image.Rect(0, 0, width, height))
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Fast background fill
	// imgdraw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.ZP, imgdraw.Src)

	w.PaintObjects(img, width, height)
	return img
}

// RenderToFile renders the world and writes the result as a PNG file.
func (w *World) RenderToFile(width, height int, bgColor color.Color, filePath string) error {
	img := w.Render(width, height, bgColor)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
