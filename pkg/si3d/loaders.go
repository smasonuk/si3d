package si3d

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"
)

// Represents a vertex with its color.
type Vertex struct {
	X, Y, Z float64
	Color   color.RGBA
}

func LoadObjectFromDXFFile(fileName string, reverse int) (*Model, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open DXF file %s: %w", fileName, err)
	}
	defer file.Close()

	obj, err := NewObjectFromDXF(file, reverse)
	if err != nil {
		return nil, fmt.Errorf("error parsing DXF file %s: %w", fileName, err)
	}

	return obj, nil
}

// NewObjectFromDXF creates a new Model by reading a simplified DXF file from
// the given reader. It returns a fully constructed object with its BSP tree
// already built, or an error if the file cannot be parsed.
func NewObjectFromDXF(reader io.Reader, reverse int) (*Model, error) {
	// create a new, empty object to populate.
	obj := NewModel()

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB max line capacity

	// Helper function to read the next line and parse it as a float64
	readFloatLine := func() (float64, error) {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF // Clean end of file
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse float value '%s': %w", scanner.Text(), err)
		}
		return val, nil
	}

	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "3DFACE") {
			continue
		}

		for i := 0; i < 3; i++ {
			if !scanner.Scan() {
				// On error, return a nil object and the error
				return nil, fmt.Errorf("unexpected end of file while parsing 3DFACE header")
			}
		}

		faceColor := color.RGBA{
			R: 100,
			G: 10,
			B: 58,
			A: 255,
		}
		aFace := NewFace(nil, faceColor, Vector3{})

		for c := 0; c < 4; c++ {
			x, err := readFloatLine()
			if err != nil {
				return nil, fmt.Errorf("error reading X coordinate for vertex %d: %w", c, err)
			}
			scanner.Scan()

			y, err := readFloatLine()
			if err != nil {
				return nil, fmt.Errorf("error reading Y coordinate for vertex %d: %w", c, err)
			}
			scanner.Scan()

			z, err := readFloatLine()
			if err != nil {
				return nil, fmt.Errorf("error reading Z coordinate for vertex %d: %w", c, err)
			}
			scanner.Scan()

			aFace.AddPoint(x, y, z)
		}

		aFace.Finished(reverse)

		//Add the face to the new object we created at the start.
		obj.faces.AddFace(aFace)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from DXF source: %w", err)
	}

	// Finalize the new object by building its BSP tree.
	obj.BuildBSP()
	obj.Compile()

	// On success, return the fully populated object and a nil error.
	return obj, nil
}

func LoadObjectFromPLYFile(fileName string, reverse int, useBsp bool) (*Model, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open PLY file %s: %w", fileName, err)
	}
	defer file.Close()

	obj, err := LoadObjectFromPLYReader(file, reverse, useBsp)
	if err != nil {
		return nil, fmt.Errorf("error parsing PLY file %s: %w", fileName, err)
	}

	return obj, nil
}

func LoadObjectFromPLYReader(reader io.Reader, reverse int, useBsp bool) (*Model, error) {
	obj := NewModel()
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB max line capacity

	var vertexCount, faceCount int
	var hasVertexColor, hasFaceColor bool
	var currentElement string

	// 1. Parse the header intelligently to detect color format
headerLoop:
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "element":
			if len(parts) == 3 {
				currentElement = parts[1] // Track if we are defining a 'vertex' or 'face'
				if parts[1] == "vertex" {
					vertexCount, _ = strconv.Atoi(parts[2])
				} else if parts[1] == "face" {
					faceCount, _ = strconv.Atoi(parts[2])
				}
			}
		case "property":
			// Check for color properties within the current element definition
			if len(parts) > 2 && (parts[2] == "red" || parts[2] == "diffuse_red") {
				if currentElement == "vertex" {
					hasVertexColor = true
				} else if currentElement == "face" {
					hasFaceColor = true
				}
			}
		case "end_header":
			break headerLoop
		}
	}

	// Read all vertices
	vertices := make([]Vertex, 0, vertexCount)
	for i := 0; i < vertexCount; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of file while reading vertices")
		}
		parts := strings.Fields(scanner.Text())

		// Conditional Vertex Parsing
		if hasVertexColor {
			if len(parts) < 6 {
				return nil, fmt.Errorf("invalid vertex-color data on line %d", i)
			}
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			z, _ := strconv.ParseFloat(parts[2], 64)
			r, _ := strconv.ParseUint(parts[3], 10, 8)
			g, _ := strconv.ParseUint(parts[4], 10, 8)
			b, _ := strconv.ParseUint(parts[5], 10, 8)
			vertices = append(vertices, Vertex{X: x, Y: y, Z: z, Color: color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}})
		} else { // No vertex color
			if len(parts) < 3 {
				return nil, fmt.Errorf("invalid vertex data on line %d", i)
			}
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			z, _ := strconv.ParseFloat(parts[2], 64)
			// Use a default white color if none is provided
			vertices = append(vertices, Vertex{X: x, Y: y, Z: z, Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}})
		}
	}

	// Read face definitions
	for i := 0; i < faceCount; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of file while reading faces")
		}
		parts := strings.Fields(scanner.Text())
		numFaceVerts, _ := strconv.Atoi(parts[0])

		var faceColor color.RGBA

		// Conditional Face Parsing
		if hasFaceColor {
			if len(parts) != numFaceVerts+1+3 {
				return nil, fmt.Errorf("invalid face-color data on line %d", i)
			}
			// Color is at the end of the line
			r, _ := strconv.ParseUint(parts[numFaceVerts+1], 10, 8)
			g, _ := strconv.ParseUint(parts[numFaceVerts+2], 10, 8)
			b, _ := strconv.ParseUint(parts[numFaceVerts+3], 10, 8)
			faceColor = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
		} else if hasVertexColor {
			// No face color, but vertices have color -> average them
			var r, g, b uint32
			for j := 1; j <= numFaceVerts; j++ {
				idx, _ := strconv.Atoi(parts[j])
				vert := vertices[idx]
				r += uint32(vert.Color.R)
				g += uint32(vert.Color.G)
				b += uint32(vert.Color.B)
			}
			faceColor = color.RGBA{
				R: uint8(r / uint32(numFaceVerts)),
				G: uint8(g / uint32(numFaceVerts)),
				B: uint8(b / uint32(numFaceVerts)),
				A: 255,
			}
		} else {
			// No color information at all, use a default color
			faceColor = color.RGBA{R: 128, G: 128, B: 128, A: 255}
		}

		aFace := NewFace(nil, faceColor, Vector3{})
		for j := 1; j <= numFaceVerts; j++ {
			idx, _ := strconv.Atoi(parts[j])
			vert := vertices[idx]
			aFace.AddPoint(vert.X, vert.Y, vert.Z)
		}

		aFace.Finished(reverse)
		obj.faces.AddFace(aFace)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from PLY source: %w", err)
	}

	if useBsp {
		obj.BuildBSP()
	}
	obj.Compile()
	return obj, nil
}
