package si3d

import (
	"image/color"
	"math"
)

// BspNode represents a node in a Binary Space Partitioning tree.
type BspNode struct {
	normal           Vector3
	Left             *BspNode
	Right            *BspNode
	colRed           uint8
	colGreen         uint8
	colBlue          uint8
	colAlpha         uint8
	facePointIndices []int
	xp               []float32
	yp               []float32
	normalIndex      int
}

// NewBspNode creates a new BSP node.
func NewBspNode(facePoints []Vector3, faceNormal Vector3, faceColor color.RGBA, pointIndices []int, normalIdx int) *BspNode {
	b := &BspNode{
		normal:           faceNormal,
		colRed:           faceColor.R,
		colGreen:         faceColor.G,
		colBlue:          faceColor.B,
		colAlpha:         faceColor.A,
		facePointIndices: pointIndices,
		normalIndex:      normalIdx,
		xp:               make([]float32, len(facePoints)),
		yp:               make([]float32, len(facePoints)),
	}
	return b
}

// set color
func (b *BspNode) SetColor(r, g, b1, a uint8) {
	b.colRed = r
	b.colGreen = g
	b.colBlue = b1
	b.colAlpha = a
}

// PaintWithoutShading paints the BSP tree without lighting effects.
func (b *BspNode) PaintWithoutShading(batcher PolygonBatcher, x, y int, transPoints []Vector3, transNormals []Vector3, linesOnly bool, screenWidth, screenHeight float32, dontDrawOutlines bool, nearPlane float64, ctx *RenderContext) {
	b.PaintWithShading(batcher, x, y, transPoints, transNormals, false, linesOnly, screenWidth, screenHeight, dontDrawOutlines, nearPlane, ctx)
}

// PaintWithShading recursively traverses the BSP tree and paints the polygons.
func (b *BspNode) PaintWithShading(batcher PolygonBatcher, x, y int, transPoints []Vector3, transNormals []Vector3, doShading bool, linesOnly bool, screenWidth, screenHeight float32, dontDrawOutlines bool, nearPlane float64, ctx *RenderContext) {
	if len(b.facePointIndices) == 0 {
		return
	}

	transformedNormal := transNormals[b.normalIndex]
	firstTransformedPoint := transPoints[b.facePointIndices[0]]

	// Determine if the polygon is facing the camera
	where := transformedNormal.X*firstTransformedPoint.X +
		transformedNormal.Y*firstTransformedPoint.Y +
		transformedNormal.Z*firstTransformedPoint.Z

	if where <= 0 { // Facing away from the camera
		if b.Left != nil {
			b.Left.PaintWithShading(batcher, x, y, transPoints, transNormals, doShading, linesOnly, screenWidth, screenHeight, dontDrawOutlines, nearPlane, ctx)
		}
		if b.Right != nil {
			b.Right.PaintWithShading(batcher, x, y, transPoints, transNormals, doShading, linesOnly, screenWidth, screenHeight, dontDrawOutlines, nearPlane, ctx)
		}
	} else { // Facing towards the camera
		if b.Right != nil {
			b.Right.PaintWithShading(batcher, x, y, transPoints, transNormals, doShading, linesOnly, screenWidth, screenHeight, dontDrawOutlines, nearPlane, ctx)
		}

		shouldReturn := b.paintPoly(batcher, x, y, transPoints, transNormals, doShading, firstTransformedPoint, transformedNormal, linesOnly, screenWidth, screenHeight, dontDrawOutlines, nearPlane, ctx)
		if shouldReturn {
			return // Z-clipping occurred, no need to paint left side
		}

		if b.Left != nil {
			b.Left.PaintWithShading(batcher, x, y, transPoints, transNormals, doShading, linesOnly, screenWidth, screenHeight, dontDrawOutlines, nearPlane, ctx)
		}
	}
}

// clipPolygonAgainstNearPlane clips a 3D polygon against the near Z plane.
func clipPolygonAgainstNearPlane(polygon []Vector3, nearPlane float64, buffer []Vector3) []Vector3 {
	if len(polygon) == 0 {
		return polygon
	}

	clippedPolygon := buffer
	startPoint := polygon[len(polygon)-1]

	for _, endPoint := range polygon {
		startInside := startPoint.Z >= nearPlane
		endInside := endPoint.Z >= nearPlane

		if endInside {
			if !startInside {
				// Edge crosses from outside to inside, calculate intersection
				intersection := intersectNearPlane(startPoint, endPoint, nearPlane)
				clippedPolygon = append(clippedPolygon, intersection)
			}
			// End point is inside, add it
			clippedPolygon = append(clippedPolygon, endPoint)
		} else if startInside {
			// Edge crosses from inside to outside, calculate intersection
			intersection := intersectNearPlane(startPoint, endPoint, nearPlane)
			clippedPolygon = append(clippedPolygon, intersection)
		}
		// If both points are outside, do nothing

		startPoint = endPoint
	}

	return clippedPolygon
}

// intersectNearPlane calculates the intersection point of a line segment with the near Z plane.
func intersectNearPlane(p1, p2 Vector3, nearPlane float64) Vector3 {
	// Using linear interpolation to find the intersection point
	// t = (nearPlaneZ - p1.z) / (p2.z - p1.z)
	deltaZ := p2.Z - p1.Z
	if math.Abs(deltaZ) < 1e-6 { // Avoid division by zero for horizontal lines
		return p1
	}
	t := (nearPlane - p1.Z) / deltaZ

	ix := p1.X + t*(p2.X-p1.X)
	iy := p1.Y + t*(p2.Y-p1.Y)

	return NewVector3(ix, iy, nearPlane)
}

// paintPoly handles Z-clipping, screen-space clipping, and drawing of a single polygon.
func (b *BspNode) paintPoly(
	batcher PolygonBatcher,
	x, y int,
	verticesInCameraSpace []Vector3,
	normalsInCameraSpace []Vector3,
	shadePoly bool,
	firstTransformedPoint Vector3,
	transformedNormal Vector3,
	linesOnly bool,
	screenWidth, screenHeight float32,
	dontDrawwOutlines bool,
	nearPlane float64,
	ctx *RenderContext,
) bool {

	initial3DPoints := ctx.Buffer3D[:0]
	for _, pointIndex := range b.facePointIndices {
		initial3DPoints = append(initial3DPoints, verticesInCameraSpace[pointIndex])
	}

	buffer := ctx.Buffer3D[len(initial3DPoints):len(initial3DPoints)]
	pointsToUse := clipPolygonAgainstNearPlane(initial3DPoints, nearPlane, buffer)

	// If clipping results in a polygon with too few vertices, don't draw it.
	if len(pointsToUse) < 3 {
		return false
	}

	ctx.BufferPoints = ctx.BufferPoints[:0]
	for _, point := range pointsToUse {
		// At this stage, point[2] (z) is guaranteed to be >= nearPlaneZ,
		// so perspective division is safe.
		z := float32(point.Z)
		ctx.BufferPoints = append(ctx.BufferPoints, Point{
			X: ConvertToScreenX(float64(screenWidth), float64(screenHeight), point.X, float64(z)),
			Y: ConvertToScreenY(float64(screenWidth), float64(screenHeight), point.Y, float64(z)),
		})
	}

	clippedPoints := batcher.ClipPolygon(ctx.BufferPoints, screenWidth, screenHeight)

	if len(clippedPoints) < 3 {
		return false
	}

	ctx.BufferFloatX = ctx.BufferFloatX[:0]
	ctx.BufferFloatY = ctx.BufferFloatY[:0]
	for _, p := range clippedPoints {
		ctx.BufferFloatX = append(ctx.BufferFloatX, p.X)
		ctx.BufferFloatY = append(ctx.BufferFloatY, p.Y)
	}

	// Copy the data to new slices so the batcher can store them
	finalScreenPointsX := make([]float32, len(ctx.BufferFloatX))
	copy(finalScreenPointsX, ctx.BufferFloatX)
	finalScreenPointsY := make([]float32, len(ctx.BufferFloatY))
	copy(finalScreenPointsY, ctx.BufferFloatY)

	polyColor := color.RGBA{R: b.colRed, G: b.colGreen, B: b.colBlue, A: b.colAlpha}
	if shadePoly {
		shadingRefPoint := verticesInCameraSpace[b.facePointIndices[0]]
		polyColor = b.calcColor(shadingRefPoint, transformedNormal, polyColor)
	}

	if !linesOnly {
		if dontDrawwOutlines {
			batcher.AddPolygon(finalScreenPointsX, finalScreenPointsY, polyColor)
		} else {
			black := color.RGBA{R: 100, G: 100, B: 100, A: 25}
			batcher.AddPolygonAndOutline(finalScreenPointsX, finalScreenPointsY, polyColor, black, 1.0)
		}
	} else {

		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		batcher.AddPolygonAndOutline(finalScreenPointsX, finalScreenPointsY, black, polyColor, 1.0)

	}

	return false
}

// calcColor calculates the color of a polygon based on simple lighting.
func (b *BspNode) calcColor(
	firstTransformedPoint Vector3,
	transformedNormal Vector3,
	polyColor color.RGBA,
) color.RGBA {
	const ambientLight = 0.65
	const spotlightConePower = 10.0
	const spotlightLightAmount = 1.0 - ambientLight

	diffuseFactor := transformedNormal.Z
	if diffuseFactor < 0 {
		diffuseFactor = 0
	}

	var spotlightFactor float64
	lenVecToPoint := GetLength2(firstTransformedPoint)

	if lenVecToPoint > 0 {
		cosAngle := firstTransformedPoint.Z / lenVecToPoint
		if cosAngle < 0 {
			cosAngle = 0
		}
		v2 := cosAngle * cosAngle
		v4 := v2 * v2
		v8 := v4 * v4
		spotlightFactor = v8 * v2
	} else {
		spotlightFactor = 1.0
	}

	spotlightBrightness := diffuseFactor * spotlightFactor * spotlightLightAmount
	finalBrightness := ambientLight + spotlightBrightness

	c := 240 - int(finalBrightness*240)
	min := 7
	r1 := clamp2(int(b.colRed)-c, min, 255)
	g1 := clamp2(int(b.colGreen)-c, min, 255)
	b1 := clamp2(int(b.colBlue)-c, min, 255)
	polyColor = color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: polyColor.A}

	return polyColor
}

// --- Helper functions ---

func clamp2(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// GetLength2 calculates the magnitude of a 3D vector.
func GetLength2(v Vector3) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func ConvertToScreenX(width, height, x, z float64) float32 {
	convFactor := width
	return float32(((convFactor * x) / z) + width/2.0)
}

func ConvertToScreenY(width, height, y, z float64) float32 {
	convFactor := width
	return float32(((convFactor * y) / z) + height/2.0)

}

// average/midpoint of the polygon in 3d space

func (b *BspNode) GetAveragePoint(transPoints []Vector3) (float64, float64, float64) {
	if len(b.facePointIndices) == 0 {
		return 0, 0, 0
	}

	sumX, sumY, sumZ := 0.0, 0.0, 0.0
	for _, index := range b.facePointIndices {
		point := transPoints[index]
		sumX += point.X
		sumY += point.Y
		sumZ += point.Z
	}

	count := float64(len(b.facePointIndices))
	return sumX / count, sumY / count, sumZ / count
}
