package si3d

import (
	"image/color"
	"math"
)

type Model struct {
	faceMesh           *FaceMesh
	normalMesh         *NormalMesh
	transFaceMesh      *FaceMesh
	transNormalMesh    *NormalMesh
	faces              *FaceStore
	root               *BspNode
	Transform          *Transform
	canPaintWithoutBSP bool
	xLength            float64
	yLength            float64
	zLength            float64
	objectDirection    Vector3
	hasObjectDirection bool
	drawLinesOnly      bool

	// for when not using BSP
	faceIndices      [][]int // Indices of the faces in the faceMesh
	normalIndices    []int
	drawAllFaces     bool // If true, draw all faces regardless of visibility
	dontDrawOutlines bool // If true, don't draw outlines of polygons
}

func (o *Model) SetDontDrawOutlines(dontDraw bool) {
	o.dontDrawOutlines = dontDraw
}

func (o *Model) GetDontDrawOutlines() bool {
	return o.dontDrawOutlines
}

func (o *Model) SetDrawAllFaces(draw bool) {
	o.drawAllFaces = draw
}

func (o *Model) SetDrawLinesOnly(only bool) {
	o.drawLinesOnly = only
}

func (o *Model) GetDrawLinesOnly() bool {
	return o.drawLinesOnly
}

func (o *Model) PaintObject(batcher *PolygonBatcher, x, y int, lightingChange bool, screenWidth, screenHeight float32, nearPlane float64, ctx *RenderContext) {
	if o.canPaintWithoutBSP {
		o.paintWithoutBSP(batcher, x, y, screenHeight, screenWidth, nearPlane, ctx)

	} else {
		if o.root != nil {
			o.root.PaintWithShading(batcher, x, y, o.transFaceMesh.Points, o.transNormalMesh.Points, lightingChange, o.drawLinesOnly, screenWidth, screenHeight,
				o.dontDrawOutlines, nearPlane, ctx)
		}
	}
}

func (o *Model) SetDirectionVector(v Vector3) {
	// normalize
	o.objectDirection = v.Copy()
	o.objectDirection = o.objectDirection.Normalize()
	o.hasObjectDirection = true

	// fmt.Println("Object direction vector set to:", o.objectDirection)
}

func (o *Model) TranslateAllPoints(x, y, z float64) {
	if o.faceMesh == nil || o.faceMesh.Points == nil {
		return
	}

	for i := range o.faceMesh.Points {
		o.faceMesh.Points[i].X += x
		o.faceMesh.Points[i].Y += y
		o.faceMesh.Points[i].Z += z
	}
}

// apply matrix to the direction vector to transform it.
func (o *Model) ApplyDirectionVector(m Matrix) Vector3 {
	if !o.hasObjectDirection {
		return Vector3{}
	}
	return m.RotateVector3(o.objectDirection)
}

// apply matrix to the direction vector to transform it.
func ApplyDirectionVector(v Vector3, m Matrix) Vector3 {
	return m.RotateVector3(v)
}

func NewModel() *Model {
	return &Model{
		transFaceMesh:   NewFaceMesh(),
		transNormalMesh: NewNormalMesh(),
		faces:           NewFaceStore(),
		Transform:       NewTransform(),
		faceIndices:     make([][]int, 0),
		normalIndices:   make([]int, 0),
	}
}

func (o *Model) Clone() *Model {
	newTransform := NewTransform()
	newTransform.Position = o.Transform.Position.Copy()
	newTransform.Rotation = o.Transform.Rotation // Value copy
	newTransform.Scale = o.Transform.Scale.Copy()

	clone := &Model{
		// shared
		faceMesh:      o.faceMesh,
		normalMesh:    o.normalMesh,
		faces:         o.faces,
		root:          o.root,
		faceIndices:   o.faceIndices,
		normalIndices: o.normalIndices,

		// instance-specific
		transFaceMesh:      o.transFaceMesh.Copy(),
		transNormalMesh:    o.transNormalMesh.Copy(),
		Transform:          newTransform,
		canPaintWithoutBSP: o.canPaintWithoutBSP,
	}
	return clone
}

func (o *Model) Compile() {
	if o.root == nil {
		o.createFaceList()
		o.canPaintWithoutBSP = true
	}

	o.faceMesh = o.transFaceMesh.Copy()
	o.normalMesh = o.transNormalMesh.Copy()

	o.CalcSize()
}

func (o *Model) BuildBSP() {
	// log.Println("Creating BSP Tree...")
	if o.faces.FaceCount() > 0 {
		o.root = o.createBspTree(o.faces, o.transFaceMesh, o.transNormalMesh)
		o.canPaintWithoutBSP = false
	}
	// log.Println("BSP Tree Created.")
}

func (o *Model) createFaceList() {
	faces, newFaces, newNormMesh := o.faces, o.transFaceMesh, o.transNormalMesh
	for i := 0; i < faces.FaceCount(); i++ {
		originalFace := faces.GetFace(i)
		newFace, ind := newFaces.AddFace(originalFace)
		newFaces.AddFace(newFace)

		normal := originalFace.GetNormal()
		_, normalIndex := newNormMesh.AddNormal(normal)

		o.faceIndices = append(o.faceIndices, ind)
		o.normalIndices = append(o.normalIndices, normalIndex)
	}
}

// Moves all points so that the 0,0,0 is the center of the object.
func (o *Model) Center() {
	if o.faceMesh == nil || len(o.faceMesh.Points) == 0 {
		return
	}

	// Calculate the bounding box of the object.
	minX, maxX := o.faceMesh.Points[0].X, o.faceMesh.Points[0].X
	minY, maxY := o.faceMesh.Points[0].Y, o.faceMesh.Points[0].Y
	minZ, maxZ := o.faceMesh.Points[0].Z, o.faceMesh.Points[0].Z
	for _, point := range o.faceMesh.Points {
		if point.X < minX {
			minX = point.X
		} else if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		} else if point.Y > maxY {
			maxY = point.Y
		}

		if point.Z < minZ {
			minZ = point.Z
		} else if point.Z > maxZ {
			maxZ = point.Z
		}
	}

	//calculate the center of the bounding box.
	centerX := (minX + maxX) / 2.0
	centerY := (minY + maxY) / 2.0
	centerZ := (minZ + maxZ) / 2.0

	// Move all points so that the center of the bounding box is at 0,0,0.
	for i := range o.faceMesh.Points {
		o.faceMesh.Points[i].X -= centerX
		o.faceMesh.Points[i].Y -= centerY
		o.faceMesh.Points[i].Z -= centerZ
	}
}

func (o *Model) GetExtents() (float64, float64, float64) {
	return o.xLength, o.yLength, o.zLength
}

func (o *Model) XLength() float64 {
	return o.xLength
}

func (o *Model) YLength() float64 {
	return o.yLength
}

func (o *Model) ZLength() float64 {
	return o.zLength
}

func (o *Model) CalcSize() {
	if o.faceMesh == nil || len(o.faceMesh.Points) == 0 {
		o.xLength = 0
		o.yLength = 0
		o.zLength = 0
		return
	}

	minX, maxX := o.faceMesh.Points[0].X, o.faceMesh.Points[0].X
	minY, maxY := o.faceMesh.Points[0].Y, o.faceMesh.Points[0].Y
	minZ, maxZ := o.faceMesh.Points[0].Z, o.faceMesh.Points[0].Z
	for _, point := range o.faceMesh.Points {
		if point.X < minX {
			minX = point.X
		} else if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		} else if point.Y > maxY {
			maxY = point.Y
		}

		if point.Z < minZ {
			minZ = point.Z
		} else if point.Z > maxZ {
			maxZ = point.Z
		}
	}
	o.xLength = maxX - minX
	o.yLength = maxY - minY
	o.zLength = maxZ - minZ
	// log.Printf("Object size: X: %.2f, Y: %.2f, Z: %.2f", o.xLength, o.yLength, o.zLength)
}

func (o *Model) ApplyMatrixTemp(aMatrix Matrix) {
	rotMatrixTemp := aMatrix.MultiplyBy(o.Transform.GetMatrix())

	// Use the new, correct method to transform the normals (rotation only).
	rotMatrixTemp.TransformNormals(o.normalMesh.Points, o.transNormalMesh.Points)

	// Use the original method to transform the vertex positions (rotation and translation).
	rotMatrixTemp.TransformObj(o.faceMesh.Points, o.transFaceMesh.Points)
}

func (o *Model) ApplyMatrixPermanent(aMatrix Matrix) {
	o.ApplyMatrixTemp(aMatrix)

	o.normalMesh.Points = o.transNormalMesh.Points // Should use Copy?
	// But NormalMesh.Points is []Vector3.
	// If ApplyMatrixTemp modifies transNormalMesh.Points in place (it does),
	// then we can just copy.
	// Wait, TransformNormals takes src, dest []Vector3.
	// dest is modified.
	// We want normalMesh.Points to be a copy of transNormalMesh.Points.
	// We can use NormalMesh.Copy() or manually copy slice.
	newPoints := make([]Vector3, len(o.transNormalMesh.Points))
	copy(newPoints, o.transNormalMesh.Points)
	o.normalMesh.Points = newPoints

	newFacePoints := make([]Vector3, len(o.transFaceMesh.Points))
	copy(newFacePoints, o.transFaceMesh.Points)
	o.faceMesh.Points = newFacePoints
}

func (o *Model) ApplyObjMatrixPermanent() {
	o.ApplyMatrixTemp(IdentMatrix())

	newPoints := make([]Vector3, len(o.transNormalMesh.Points))
	copy(newPoints, o.transNormalMesh.Points)
	o.normalMesh.Points = newPoints

	newFacePoints := make([]Vector3, len(o.transFaceMesh.Points))
	copy(newFacePoints, o.transFaceMesh.Points)
	o.faceMesh.Points = newFacePoints

	o.Transform = NewTransform()
}

func (o *Model) createBspTree(faces *FaceStore, newFaces *FaceMesh, newNormMesh *NormalMesh) *BspNode {
	if faces.FaceCount() == 0 {
		return nil
	}

	parentFace := o.choosePlane(faces)
	originalNormal, normalIndex := newNormMesh.AddNormal(parentFace.GetNormal())
	parentFace.SetNormal(NewVector3(originalNormal.X, originalNormal.Y, originalNormal.Z))
	newFace, parentIndices := newFaces.AddFace(parentFace)
	parent := NewBspNode(newFace.Points, newFace.GetNormal(), newFace.Col, parentIndices, normalIndex)
	pPlane := NewPlane(newFace, newFace.GetNormal())

	// Create two new lists to hold the faces that fall on either side of the plane.
	fvLeft := NewFaceStore()
	fvRight := NewFaceStore()

	// Partition the *remaining* faces against the splitting plane.
	for a := 0; a < faces.FaceCount(); a++ {
		currentFace := faces.GetFace(a)
		if pPlane.FaceIntersect(currentFace) {
			// If the face is split by the plane...
			split := pPlane.SplitFace(currentFace)
			if split == nil {
				continue
			}
			for _, facePart := range split {
				if facePart != nil && len(facePart.Points) > 0 {
					if pPlane.Where(facePart) <= 0 {
						f1 := NewFace(facePart.Points, currentFace.Col, currentFace.GetNormal())
						fvLeft.AddFace(f1)
					} else {
						f2 := NewFace(facePart.Points, currentFace.Col, currentFace.GetNormal())
						fvRight.AddFace(f2)
					}
				}
			}
		} else {
			// If the face is entirely on one side...
			w := pPlane.Where(currentFace)
			if w <= 0 {
				fvLeft.AddFace(currentFace)
			} else {
				fvRight.AddFace(currentFace)
			}
		}
	}

	// Build the left and right sub-trees from the new lists.
	if fvLeft.FaceCount() > 0 {
		parent.Left = o.createBspTree(fvLeft, newFaces, newNormMesh)
	}
	if fvRight.FaceCount() > 0 {
		parent.Right = o.createBspTree(fvRight, newFaces, newNormMesh)
	}

	return parent
}

func (o *Model) choosePlane(fs *FaceStore) *Face {
	const maxCandidates = 100
	leastFace, leastFaceTotal := 0, fs.FaceCount()

	numFaces := fs.FaceCount()
	if numFaces == 0 {
		return nil
	}

	step := 1
	if numFaces > maxCandidates {
		step = numFaces / maxCandidates
	}

	for chosenCandidate := 0; chosenCandidate < numFaces; chosenCandidate += step {
		total := 0
		p := fs.GetFace(chosenCandidate).GetPlane()
		for i := 0; i < fs.FaceCount(); i++ {
			if i == chosenCandidate {
				continue
			}
			f := fs.GetFace(i)
			if p.FaceIntersect(f) {
				total++
			}
		}
		if total < leastFaceTotal {
			leastFaceTotal = total
			leastFace = chosenCandidate
			if total == 0 {
				break
			}
		}
	}
	return fs.RemoveFaceAt(leastFace)
}

func (o *Model) GetPosition() Vector3 {
	return o.Transform.Position
}

func (o *Model) SetPosition(x, y, z float64) {
	o.Transform.Position = NewVector3(x, y, z)
}

func (o *Model) RollObject(directionOfRoll float64, amountOfMovement float64) {
	// Simplified rolling: update X (Pitch) and Z (Roll) rotation based on direction
	o.Transform.RotateLocal(NewVector3(1, 0, 0), -amountOfMovement*math.Cos(directionOfRoll))
	o.Transform.Rotate(NewVector3(0, 0, 1), amountOfMovement*math.Sin(directionOfRoll))
}

func (o *Model) RotateY(amountOfMovementInRads float64) {
	o.Transform.Rotate(NewVector3(0, 1, 0), amountOfMovementInRads)
}

func ApplyRollObjectMatrix(directionOfRoll float64, amountOfMovement float64, existingMatrix Matrix) Matrix {
	rotMatrixY := NewRotationMatrix(ROTY, -directionOfRoll)
	rotMatrixYBack := NewRotationMatrix(ROTY, directionOfRoll)
	rotMatrixX := NewRotationMatrix(ROTX, -amountOfMovement)
	all := rotMatrixY.MultiplyBy(rotMatrixX).MultiplyBy(rotMatrixYBack)

	return all.MultiplyBy(existingMatrix)
}

func ApplyRotateYObjectMatrix(amountOfMovementInRads float64, existingMatrix Matrix) Matrix {
	rotMatrix := NewRotationMatrix(ROTY, amountOfMovementInRads)
	return rotMatrix.MultiplyBy(existingMatrix)
}

func (o *Model) AddFacesFromObject(other *Model) {
	if other == nil || other.faces == nil || other.faces.FaceCount() == 0 {
		return
	}

	for i := 0; i < other.faces.FaceCount(); i++ {
		face := other.faces.GetFace(i)
		newFace := face.Copy()

		o.faces.AddFace(newFace)
	}
}

// get face indices
func (o *Model) GetFaceIndices() [][]int {
	return o.faceIndices
}

// get transFaceMesh
func (o *Model) GetTransFaceMesh() *FaceMesh {
	return o.transFaceMesh
}

func (o *Model) paintWithoutBSP(batcher *PolygonBatcher, x, y int, screenHeight, screenWidth float32, nearPlane float64, ctx *RenderContext) {
	for i := 0; i < len(o.faceIndices); i++ {
		faceIndices := o.faceIndices[i]
		normalIndex := o.normalIndices[i]
		facePointsInCameraSpace := ctx.Buffer3D[:0]
		for _, index := range faceIndices {
			point := o.transFaceMesh.Points[index]
			facePointsInCameraSpace = append(facePointsInCameraSpace, point)
		}
		face := o.faces.faces[i]

		normal := o.transNormalMesh.Points[normalIndex]

		o.paintFace(batcher, x, y, facePointsInCameraSpace, normal, screenWidth, screenHeight, face, nearPlane, ctx)
	}
}

func (o *Model) BspNodesIntersectingLine(startLine, endLine Vector3) *BspNode {
	nodes := make([]*BspNode, 0, 5)
	points := make([]Vector3, 0, 10)

	var traverseNodes func(node *BspNode)
	traverseNodes = func(node *BspNode) {
		if node == nil {
			return
		}

		facePointsInCameraSpace := points[:]
		for _, index := range node.facePointIndices {
			point := o.transFaceMesh.Points[index]
			facePointsInCameraSpace = append(facePointsInCameraSpace, point)
		}
		intersects := LineIntersectsPolygon(startLine, endLine, facePointsInCameraSpace)
		if intersects {
			nodes = append(nodes, node)
		}

		traverseNodes(node.Left)
		traverseNodes(node.Right)
	}

	traverseNodes(o.root)

	if len(nodes) == 0 {
		return nil
	}

	// loop through the nodes and find the one closest to startLine point
	closestNode := nodes[0]
	distance := math.MaxFloat64
	for _, node := range nodes {
		nodex, nodey, nodez := node.GetAveragePoint(o.transFaceMesh.Points)
		dx := nodex - startLine.X
		dy := nodey - startLine.Y
		dz := nodez - startLine.Z
		d := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if d < distance {
			distance = d
			closestNode = node
		}
	}

	return closestNode
}

func (o *Model) paintFace(batcher *PolygonBatcher, x, y int, points []Vector3, normal Vector3, screenWidth, screenHeight float32, face *Face, nearPlane float64, ctx *RenderContext) {

	firstPoint := points[0]
	where := 1.0
	if !o.drawAllFaces {
		where = normal.X*firstPoint.X +
			normal.Y*firstPoint.Y +
			normal.Z*firstPoint.Z
	}

	if where > 0 { // Facing the camera
		o.paintFace2(batcher,
			x,
			y,
			firstPoint,
			points,
			normal,
			screenWidth,
			screenHeight,
			face,
			nearPlane,
			ctx,
		)
	}

}

func (o *Model) paintFace2(
	batcher *PolygonBatcher,
	x, y int,
	firstTransformedPoint Vector3,
	initial3DPoints []Vector3,
	transformedNormal Vector3,
	screenWidth, screenHeight float32,
	face *Face,
	nearPlane float64,
	ctx *RenderContext,
) bool {

	buffer := ctx.Buffer3D[len(initial3DPoints):len(initial3DPoints)]
	pointsToUse := clipPolygonAgainstNearPlane(initial3DPoints, nearPlane, buffer)
	if len(pointsToUse) < 3 {
		return false
	}

	ctx.BufferPoints = ctx.BufferPoints[:0]
	for _, point := range pointsToUse {
		// cf := float64(screenWidth)
		// At this stage, point[2] (z) is guaranteed to be >= nearPlaneZ,
		// so perspective division is safe.
		z := float32(point.Z)
		ctx.BufferPoints = append(ctx.BufferPoints, Point{
			X: ConvertToScreenX(float64(screenWidth), float64(screenHeight), point.X, float64(z)),
			Y: ConvertToScreenY(float64(screenWidth), float64(screenHeight), point.Y, float64(z)),
		})
	}

	// clip to screen
	clippedPoints := batcher.clipPolygon(ctx.BufferPoints, screenWidth, screenHeight)
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

	col := face.Col
	polyColor := col
	if true {
		shadingRefPoint := firstTransformedPoint
		polyColor = getColor(shadingRefPoint, transformedNormal, polyColor)
	}

	if !o.drawLinesOnly {
		black := color.RGBA{R: 50, G: 50, B: 50, A: 25}
		batcher.AddPolygonAndOutline(finalScreenPointsX, finalScreenPointsY, polyColor, black, 1.0)

	} else {
		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		batcher.AddPolygonAndOutline(finalScreenPointsX, finalScreenPointsY, black, polyColor, 1.0)
	}

	return false
}

// GetColor calculates the color of a polygon based on simple lighting.
func getColor(
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
	r1 := clamp2(int(polyColor.R)-c, min, 255)
	g1 := clamp2(int(polyColor.G)-c, min, 255)
	b1 := clamp2(int(polyColor.B)-c, min, 255)
	polyColor = color.RGBA{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: polyColor.A}

	return polyColor
}
func (o *Model) ScaleAllPoints(scale float64) {
	if o.faceMesh == nil || o.faceMesh.Points == nil {
		return
	}

	for i := range o.faceMesh.Points {
		o.faceMesh.Points[i].X *= scale
		o.faceMesh.Points[i].Y *= scale
		o.faceMesh.Points[i].Z *= scale
	}
}

// A small epsilon value for floating-point comparisons to avoid precision errors.
const epsilon = 1e-6

// LineIntersectsPolygon determines if a line segment intersects with a 3D polygon.
// The polygon is assumed to be planar and convex.
func LineIntersectsPolygon(lineStart, lineEnd Vector3, polygonPoints []Vector3) bool {
	if len(polygonPoints) < 3 {
		// A polygon must have at least 3 vertices.
		return false
	}

	// 1. Define the plane of the polygon using the first three points.
	// We need a point on the plane (p0) and the plane's normal vector.
	p0 := polygonPoints[0]
	p1 := polygonPoints[1]
	p2 := polygonPoints[2]

	// Create two vectors on the plane.
	v1 := Subtract(p1, p0)
	v2 := Subtract(p2, p0)

	// The plane normal is the cross product of the two vectors.
	planeNormal := Cross(v1, v2)

	// 2. Calculate the intersection of the line and the plane.
	// The line is represented as P = lineStart + t * lineDir
	lineDir := Subtract(lineEnd, lineStart)

	// Check if the line is parallel to the plane.
	dotNormalDir := Dot(planeNormal, lineDir)
	if math.Abs(dotNormalDir) < epsilon {
		return false // Line is parallel, no intersection.
	}

	// Calculate the 't' parameter for the line equation.
	w := Subtract(lineStart, p0)
	t := -Dot(planeNormal, w) / dotNormalDir

	// 3. Check if the intersection point is within the line segment.
	// If t is not between 0 and 1, the intersection is outside the segment.
	if t < 0.0-epsilon || t > 1.0+epsilon {
		return false
	}

	// 4. Calculate the actual intersection point.
	intersectionPoint := NewVector3(
		lineStart.X+t*lineDir.X,
		lineStart.Y+t*lineDir.Y,
		lineStart.Z+t*lineDir.Z,
	)

	// 5. Check if the intersection point is inside the polygon.
	return isPointInPolygon(intersectionPoint, polygonPoints, planeNormal)
}

// isPointInPolygon checks if a 3D point (known to be on the polygon's plane)
// is inside the polygon's boundaries using a 2D projection and ray casting.
func isPointInPolygon(point Vector3, polygonPoints []Vector3, normal Vector3) bool {
	// Project the 3D polygon and the point to a 2D plane.
	// We choose the plane based on the largest component of the normal vector
	// to avoid a degenerate projection (i.e., the polygon projecting to a line).
	absX := math.Abs(normal.X)
	absY := math.Abs(normal.Y)
	absZ := math.Abs(normal.Z)

	var u, v int // Indices for the 2D coordinates (0=X, 1=Y, 2=Z)

	if absX > absY && absX > absZ {
		// Project to YZ plane (discarding X)
		u, v = 1, 2
	} else if absY > absX && absY > absZ {
		// Project to XZ plane (discarding Y)
		u, v = 0, 2
	} else {
		// Project to XY plane (discarding Z)
		u, v = 0, 1
	}

	// Get the 2D coordinates of the intersection point.
	point2D := []float64{getCoord(point, u), getCoord(point, v)}

	// Apply the Ray Casting algorithm in 2D.
	intersections := 0
	numVertices := len(polygonPoints)
	for i := 0; i < numVertices; i++ {
		p1_3D := polygonPoints[i]
		p2_3D := polygonPoints[(i+1)%numVertices]

		// Get 2D coordinates of the edge's vertices.
		p1_2D := []float64{getCoord(p1_3D, u), getCoord(p1_3D, v)}
		p2_2D := []float64{getCoord(p2_3D, u), getCoord(p2_3D, v)}

		// Check if the horizontal ray from the point intersects with the edge.
		if (p1_2D[1] > point2D[1]) != (p2_2D[1] > point2D[1]) {
			// Calculate the x-intersection of the line.
			x_intersection := (p2_2D[0]-p1_2D[0])*(point2D[1]-p1_2D[1])/(p2_2D[1]-p1_2D[1]) + p1_2D[0]
			if point2D[0] < x_intersection {
				intersections++
			}
		}
	}

	// If the number of intersections is odd, the point is inside the polygon.
	return intersections%2 == 1
}

// getCoord is a small helper to get a coordinate by its index (0=X, 1=Y, 2=Z).
func getCoord(p Vector3, index int) float64 {
	switch index {
	case 0:
		return p.X
	case 1:
		return p.Y
	case 2:
		return p.Z
	}
	return 0
}
