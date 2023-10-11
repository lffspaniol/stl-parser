package stl

// Triangle is a triangle in 3D space.
type Triangle struct {
	// Normal is the normal vector of the triangle.
	Normal [3]float32

	// Vertices are the vertices of the triangle.
	Vertices [3][3]float32

	// Attr is the attribute byte count.
	Attr uint16
}
