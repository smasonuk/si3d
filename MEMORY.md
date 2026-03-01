# MEMORY.md - si3d (Software 3D Rendering Engine)

## 📌 Project Overview
`si3d` is a custom, CPU-based software 3D rendering engine written in Go. It performs 3D math, transformations, depth sorting, clipping, and rasterization entirely in software without relying on hardware-accelerated graphics APIs like OpenGL or Vulkan. It outputs rendered scenes to standard Go `image.RGBA` buffers (which can be saved to disk, e.g., as PNGs).

## 🏗️ Architecture & Core Components

### 1. Math & Transformations
* **`Vector3`**: Standard 3D vector representation with helper methods (Add, Subtract, Cross, Dot, Normalize).
* **`Matrix`**: Custom fixed-size 4x4 matrix (`[4][4]float64`) optimized for cache locality. Handles translation, scaling, and rotations.
* **`Transform`**: Represents an object's position, scale, and rotation. Uses **Quaternions** (`github.com/go-gl/mathgl/mgl64`) to handle rotations and avoid gimbal lock.

### 2. Scene Graph
* **`World` (`world.go`)**: The root container for the scene. Holds cameras and entities. It manages the high-level rendering loop, including sorting objects by distance to the camera for correct draw order (painters algorithm for objects).
* **`Camera` (`camera.go`)**: Defines the viewpoint. Supports `LookAt` targeting and uses Quaternions for internal rotation tracking. Maintains a near clipping plane.
* **`Entity`**: A spatial instance of a `Model` placed in the `World` at a specific `X, Y, Z` coordinate.

### 3. Geometry & Meshes
* **`Mesh` / `FaceMesh` / `NormalMesh`**: Core data structures storing vertices and normals. Uses a map-based index (`pointIndex`) for $O(1)$ duplicate vertex lookups during mesh construction.
* **`Face` / `Plane`**: Represents polygonal surfaces (triangles/quads). `Plane` calculates intersection math necessary for BSP splitting and clipping.
* **`Model` (`model.go`)**: The primary 3D object representation. Contains transformed and untransformed meshes. It supports two rendering paths:
    1.  **BSP Tree Rendering**: For complex/concave objects, it generates a Binary Space Partitioning (`BspNode`) tree to guarantee correct back-to-front rendering of overlapping faces. 
    2.  **Simple Rendering**: For simple objects or wireframes where strict face depth-sorting isn't required.

### 4. Rendering Pipeline

1.  **Transformation**: Vertices are transformed from Local Space -> World Space -> Camera Space.
2.  **Z-Sorting**: The `World` sorts background/foreground objects based on their distance from the camera.
3.  **BSP Traversal**: For BSP-enabled models, the tree is traversed to draw polygons back-to-front relative to the camera position.
4.  **3D Near-Plane Clipping**: Polygons are clipped against the camera's near Z-plane (`clipPolygonAgainstNearPlane`) to prevent behind-camera vertices from mirroring or causing divide-by-zero errors.
5.  **Projection**: 3D coordinates are projected onto the 2D screen using perspective division (`ConvertToScreenX`, `ConvertToScreenY`).
6.  **2D Screen Clipping**: The `DefaultBatcher` applies the **Sutherland-Hodgman algorithm** to clip 2D polygons strictly to the screen dimensions.
7.  **Rasterization**: The 2D clipped polygons are batched and drawn to an `image.RGBA` using the `github.com/fogleman/gg` 2D rendering library. Features simple flat shading/lighting based on face normals and a simulated spotlight.

### 5. Generators, Loaders & Exporters
* **Solids Generator (`model_creators_solids.go`)**: Procedural generation of primitive shapes: Cubes, Rectangles, Cylinders, Rings, Spheres (Icosahedron/UVSphere subdivision), and Subdivided Planes.
* **Heightmaps**: Supports procedural terrain generation using Perlin Noise (`github.com/aquilax/go-perlin`).
* **Loaders (`loaders.go`)**: Can parse `.PLY` (Polygon File Format) and `.DXF` files into `Model` structs.
* **Exporters (`exporters.go`)**: Can export a built `Model` back out to `.DXF` or `.PLY` (with face colors).

## 🛠️ Tech Stack & Dependencies
* **Language**: Go
* **Dependencies**: 
    * `github.com/go-gl/mathgl/mgl64` (Quaternion and Matrix math)
    * `github.com/fogleman/gg` (2D rasterization/drawing context)
    * `github.com/aquilax/go-perlin` (Procedural terrain generation)

## 📝 Current State & TODOs
* **Rendering Status**: Fully functional wireframe and flat-shaded polygon rendering. E2E tests (`render_e2e_test.go`) use golden image regression testing to verify output stability.
* **Known TODOs (from `README.md`)**:
    * Implement clipping on the full viewing frustum (currently, only near-plane Z-clipping and 2D screen-edge clipping are implemented; left/right/top/bottom 3D frustum clipping is missing).