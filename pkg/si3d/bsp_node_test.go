package si3d

import (
	"image/color"
	"math"
	"testing"
)

func TestNewBspNode(t *testing.T) {
	facePoints := []Vector3{NewVector3(0, 0, 0), NewVector3(1, 0, 0), NewVector3(0, 1, 0)}
	faceNormal := NewVector3(0, 0, 1)
	col := color.RGBA{255, 0, 0, 255}
	pointIndices := []int{0, 1, 2}

	node := NewBspNode(facePoints, faceNormal, col, pointIndices, 0)

	if node == nil {
		t.Fatal("NewBspNode returned nil")
	}
	if node.colRed != 255 {
		t.Errorf("NewBspNode color incorrect")
	}
	if len(node.facePointIndices) != 3 {
		t.Errorf("NewBspNode indices incorrect")
	}
}

func TestIntersectNearPlane(t *testing.T) {
	// nearPlaneZ is 10 (from bsp_node.go)
	nearPlaneZ := 10.0

	// Line crossing the near plane
	p1 := NewVector3(0, 0, 0)
	p2 := NewVector3(0, 0, 20)

	intersection := intersectNearPlane(p1, p2, nearPlaneZ)

	if intersection.Z != nearPlaneZ {
		t.Errorf("intersectNearPlane Z expected %f, got %f", nearPlaneZ, intersection.Z)
	}
	// Should be halfway because 10 is halfway between 0 and 20
	if intersection.X != 0 || intersection.Y != 0 {
		t.Errorf("intersectNearPlane X,Y expected 0,0, got %f,%f", intersection.X, intersection.Y)
	}

	// Another case
	p3 := NewVector3(0, 0, 5)
	p4 := NewVector3(10, 0, 15)
	// Plane at 10.
	// t = (10 - 5) / (15 - 5) = 5 / 10 = 0.5
	// x = 0 + 0.5 * (10 - 0) = 5

	intersection2 := intersectNearPlane(p3, p4, nearPlaneZ)
	if intersection2.Z != nearPlaneZ {
		t.Errorf("intersectNearPlane2 Z expected %f, got %f", nearPlaneZ, intersection2.Z)
	}
	if math.Abs(intersection2.X-5.0) > 1e-6 {
		t.Errorf("intersectNearPlane2 X expected 5, got %f", intersection2.X)
	}
}

func TestClipPolygonAgainstNearPlane(t *testing.T) {
	// nearPlaneZ = 10
	nearPlaneZ := 10.0

	// Triangle completely behind near plane (Z < 10)
	poly1 := []Vector3{
		NewVector3(0, 0, 0),
		NewVector3(1, 0, 0),
		NewVector3(0, 1, 0),
	}
	buffer1 := make([]Vector3, 0, 10)
	clipped1 := clipPolygonAgainstNearPlane(poly1, nearPlaneZ, buffer1)
	if len(clipped1) != 0 {
		t.Errorf("Expected clipped1 to be empty, got %d points", len(clipped1))
	}

	// Triangle completely in front of near plane (Z >= 10)
	poly2 := []Vector3{
		NewVector3(0, 0, 11),
		NewVector3(1, 0, 11),
		NewVector3(0, 1, 11),
	}
	buffer2 := make([]Vector3, 0, 10)
	clipped2 := clipPolygonAgainstNearPlane(poly2, nearPlaneZ, buffer2)
	if len(clipped2) != 3 {
		t.Errorf("Expected clipped2 to have 3 points, got %d", len(clipped2))
	}

	// Triangle crossing near plane
	poly3 := []Vector3{
		NewVector3(0, 0, 5),   // Outside
		NewVector3(0, 0, 15),  // Inside
		NewVector3(10, 0, 15), // Inside
	}
	// Should clip.
	// Edge 1: 5 -> 15. Enters. Intersection at 10. Add Intersection, Add End. (2 points)
	// Edge 2: 15 -> 15. Inside -> Inside. Add End. (1 point)
	// Edge 3: 15 -> 5. Exits. Intersection at 10. Add Intersection. (1 point)
	// Total 4 points.

	buffer3 := make([]Vector3, 0, 10)
	clipped3 := clipPolygonAgainstNearPlane(poly3, nearPlaneZ, buffer3)
	if len(clipped3) != 4 {
		t.Errorf("Expected clipped3 to have 4 points, got %d", len(clipped3))
	}
}

func TestClipPolygon(t *testing.T) {
	// Clip against screen (0,0) to (width, height)
	width := float32(100)
	height := float32(100)

	batcher := NewDefaultBatcher(1)

	// Square inside screen
	poly1 := []Point{
		{10, 10},
		{90, 10},
		{90, 90},
		{10, 90},
	}

	clipped1 := batcher.ClipPolygon(poly1, width, height)
	if len(clipped1) != 4 {
		t.Errorf("Expected clipped1 to have 4 points, got %d", len(clipped1))
	}

	// Square larger than screen
	poly2 := []Point{
		{-10, -10},
		{110, -10},
		{110, 110},
		{-10, 110},
	}

	// Should be clipped to screen rect (4 corners)
	clipped2 := batcher.ClipPolygon(poly2, width, height)
	if len(clipped2) != 4 {
		t.Errorf("Expected clipped2 to have 4 points, got %d", len(clipped2))
	}
	// Check approximate coords (might have +1 hack in code)
	// The code adds 1 to width/height

	// Triangle with one point outside
	poly3 := []Point{
		{50, 50},
		{150, 50}, // Outside right
		{50, 150}, // Outside bottom
	}

	clipped3 := batcher.ClipPolygon(poly3, width, height)
	// Should result in a polygon with more vertices potentially, but at least clipped.
	// 50,50 is in.
	// 150,50 is out. Clip against right edge. Intersection at 100,50.
	// ...
	if len(clipped3) < 3 {
		t.Errorf("Expected clipped3 to be a valid polygon, got %d points", len(clipped3))
	}
}

func TestBspNode_GetAveragePoint(t *testing.T) {
	node := &BspNode{
		facePointIndices: []int{0, 1},
	}
	points := []Vector3{
		NewVector3(0, 0, 0),
		NewVector3(10, 10, 10),
	}

	x, y, z := node.GetAveragePoint(points)
	if x != 5 || y != 5 || z != 5 {
		t.Errorf("GetAveragePoint expected 5,5,5, got %f,%f,%f", x, y, z)
	}
}

func TestClipPolygonAgainstNearPlane_InFront(t *testing.T) {
	nearPlaneZ := 10.0
	// Polygon entirely in front of nearPlaneZ (Z >= 10)
	poly := []Vector3{
		NewVector3(0, 0, 15),
		NewVector3(10, 0, 15),
		NewVector3(0, 10, 15),
	}
	buffer := make([]Vector3, 0, 10)
	clipped := clipPolygonAgainstNearPlane(poly, nearPlaneZ, buffer)

	if len(clipped) != 3 {
		t.Fatalf("Expected 3 points, got %d", len(clipped))
	}
	for i, p := range clipped {
		if p != poly[i] {
			t.Errorf("Point mismatch at index %d: expected %v, got %v", i, poly[i], p)
		}
	}
}

func TestClipPolygonAgainstNearPlane_Behind(t *testing.T) {
	nearPlaneZ := 10.0
	// Polygon entirely behind nearPlaneZ (Z < 10)
	poly := []Vector3{
		NewVector3(0, 0, 5),
		NewVector3(10, 0, 5),
		NewVector3(0, 10, 5),
	}
	buffer := make([]Vector3, 0, 10)
	clipped := clipPolygonAgainstNearPlane(poly, nearPlaneZ, buffer)

	if len(clipped) != 0 {
		t.Errorf("Expected empty slice, got %d points", len(clipped))
	}
}

func TestClipPolygonAgainstNearPlane_Straddling(t *testing.T) {
	nearPlaneZ := 10.0
	// Polygon straddling nearPlaneZ (10)
	// Point 1: (0, 0, 5)  -> Outside
	// Point 2: (0, 0, 15) -> Inside
	// Point 3: (10, 0, 15) -> Inside
	poly := []Vector3{
		NewVector3(0, 0, 5),
		NewVector3(0, 0, 15),
		NewVector3(10, 0, 15),
	}

	buffer := make([]Vector3, 0, 10)
	clipped := clipPolygonAgainstNearPlane(poly, nearPlaneZ, buffer)

	// Expected order (based on Sutherland-Hodgman implementation traversing len-1 to 0):
	// 1. Edge P3(15)->P1(5): Exits. Intersection at Z=10. Add Intersection (5, 0, 10).
	// 2. Edge P1(5)->P2(15): Enters. Intersection at Z=10. Add Intersection (0, 0, 10). Add End P2(0, 0, 15).
	// 3. Edge P2(15)->P3(15): Inside. Add End P3(10, 0, 15).
	// Result: [(5, 0, 10), (0, 0, 10), (0, 0, 15), (10, 0, 15)]

	if len(clipped) != 4 {
		t.Fatalf("Expected 4 points, got %d", len(clipped))
	}

	// Verify Z coordinates
	for i, p := range clipped {
		if p.Z < nearPlaneZ {
			t.Errorf("Point %d has Z < %f: %f", i, nearPlaneZ, p.Z)
		}
	}

	// Check specific intersection points

	// Point 0: Intersection (5,0,10)
	if clipped[0].Z != nearPlaneZ {
		t.Errorf("Expected clipped[0].Z to be %f, got %f", nearPlaneZ, clipped[0].Z)
	}
	if math.Abs(clipped[0].X-5.0) > 1e-6 {
		t.Errorf("Expected clipped[0].X to be 5, got %f", clipped[0].X)
	}

	// Point 1: Intersection (0,0,10)
	if clipped[1].Z != nearPlaneZ {
		t.Errorf("Expected clipped[1].Z to be %f, got %f", nearPlaneZ, clipped[1].Z)
	}
	if clipped[1].X != 0 || clipped[1].Y != 0 {
		t.Errorf("Expected clipped[1] to be (0,0,10), got %v", clipped[1])
	}
}
