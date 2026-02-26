package si3d

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// Point represents a 2D point in screen space.
type Point struct {
	X, Y float32
}

// polygonCommand stores one deferred polygon draw operation.
type polygonCommand struct {
	xp        []float32
	yp        []float32
	fillClr   color.RGBA
	strokeClr color.RGBA
	strokeW   float32
	hasStroke bool
	hasFill   bool
}

type PolygonBatcher struct {
	commands []polygonCommand
	clipBufA []Point
	clipBufB []Point
}

func NewPolygonBatcher(initialCap int) *PolygonBatcher {
	return &PolygonBatcher{
		commands: make([]polygonCommand, 0, initialCap),
		clipBufA: make([]Point, 0, 20),
		clipBufB: make([]Point, 0, 20),
	}
}

// AddPolygon adds a single polygon's geometry to the batch.
func (b *PolygonBatcher) AddPolygon(xp, yp []float32, clr color.RGBA) {
	if len(xp) < 3 {
		return
	}

	cmd := polygonCommand{
		xp:      xp,
		yp:      yp,
		fillClr: clr,
		hasFill: true,
	}
	b.commands = append(b.commands, cmd)
}

func (b *PolygonBatcher) AddPolygonOutline(xp, yp []float32, strokeWidth float32, clr color.RGBA) {
	if len(xp) < 2 {
		return
	}

	cmd := polygonCommand{
		xp:        xp,
		yp:        yp,
		strokeClr: clr,
		strokeW:   strokeWidth,
		hasStroke: true,
	}
	b.commands = append(b.commands, cmd)
}

// AddPolygonWithOutline adds a filled polygon and its outline to the batch.
func (b *PolygonBatcher) AddPolygonAndOutline(xp, yp []float32, fillClr, strokeClr color.RGBA, strokeWidth float32) {
	if len(xp) < 3 {
		return // Need at least 3 vertices for a polygon.
	}

	cmd := polygonCommand{
		xp:        xp,
		yp:        yp,
		fillClr:   fillClr,
		strokeClr: strokeClr,
		strokeW:   strokeWidth,
		hasFill:   true,
		hasStroke: true,
	}
	b.commands = append(b.commands, cmd)
}

// Draw sends the entire batch of polygons to be drawn on the target image.
func (b *PolygonBatcher) Draw(target *image.RGBA) {
	if len(b.commands) == 0 {
		return
	}

	dc := gg.NewContextForRGBA(target)

	for _, cmd := range b.commands {
		if len(cmd.xp) == 0 {
			continue
		}

		dc.NewSubPath()
		dc.MoveTo(float64(cmd.xp[0]), float64(cmd.yp[0]))
		for i := 1; i < len(cmd.xp); i++ {
			dc.LineTo(float64(cmd.xp[i]), float64(cmd.yp[i]))
		}
		dc.ClosePath()

		if cmd.hasFill {
			dc.SetColor(cmd.fillClr)
			if cmd.hasStroke {
				dc.FillPreserve()
			} else {
				dc.Fill()
			}
		}

		if cmd.hasStroke {
			dc.SetColor(cmd.strokeClr)
			dc.SetLineWidth(float64(cmd.strokeW))
			dc.Stroke()
		}
	}

	// Reset the commands slice for the next frame, but keep the allocated memory.
	b.commands = b.commands[:0]
}

// clipPolygon applies the Sutherland-Hodgman algorithm to clip a polygon against the screen boundaries.
func (b *PolygonBatcher) clipPolygon(subjectPolygon []Point, screenWidth, screenHeight float32) []Point {
	// a bit of a hack so the edges of the new polygon are not exactly on the screen edges.
	screenWidth = screenWidth + 1
	screenHeight = screenHeight + 1
	// Clip against the 4 screen edges sequentially.
	clipped := b.clipAgainstEdge(subjectPolygon, b.clipBufA, func(p Point) bool { return p.X >= 0 }, func(s, e Point) Point { // Left edge
		if e.X == s.X {
			return Point{X: 0, Y: s.Y}
		}
		return Point{X: 0, Y: s.Y + (e.Y-s.Y)*(0-s.X)/(e.X-s.X)}
	})
	clipped = b.clipAgainstEdge(clipped, b.clipBufB, func(p Point) bool { return p.X <= screenWidth }, func(s, e Point) Point { // Right edge
		if e.X == s.X {
			return Point{X: screenWidth, Y: s.Y}
		}
		return Point{X: screenWidth, Y: s.Y + (e.Y-s.Y)*(screenWidth-s.X)/(e.X-s.X)}
	})
	clipped = b.clipAgainstEdge(clipped, b.clipBufA, func(p Point) bool { return p.Y >= 0 }, func(s, e Point) Point { // Top edge
		if e.Y == s.Y {
			return Point{X: s.X, Y: 0}
		}
		return Point{X: s.X + (e.X-s.X)*(0-s.Y)/(e.Y-s.Y), Y: 0}
	})
	clipped = b.clipAgainstEdge(clipped, b.clipBufB, func(p Point) bool { return p.Y <= screenHeight }, func(s, e Point) Point { // Bottom edge
		if e.Y == s.Y {
			return Point{X: s.X, Y: screenHeight}
		}
		return Point{X: s.X + (e.X-s.X)*(screenHeight-s.Y)/(e.Y-s.Y), Y: screenHeight}
	})

	return clipped
}

// clipAgainstEdge clips a polygon against a single arbitrary edge.
func (b *PolygonBatcher) clipAgainstEdge(subjectPolygon []Point, outputBuffer []Point, inside func(Point) bool, intersection func(Point, Point) Point) []Point {
	if len(subjectPolygon) == 0 {
		return subjectPolygon
	}

	outputBuffer = outputBuffer[:0]
	s := subjectPolygon[len(subjectPolygon)-1]

	for _, e := range subjectPolygon {
		sInside := inside(s)
		eInside := inside(e)

		if eInside {
			if !sInside {
				outputBuffer = append(outputBuffer, intersection(s, e))
			}
			outputBuffer = append(outputBuffer, e)
		} else if sInside {
			outputBuffer = append(outputBuffer, intersection(s, e))
		}

		s = e
	}

	return outputBuffer
}
