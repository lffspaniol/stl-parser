package stl

import (
	"stl-parser/internal/file"
)

var name string
var triangles []Triangle

func CountTriangles(filepath string) (string, int, error) {
	scanner, close, err := file.Reader(filepath)
	if err != nil {
		return "", -1, err
	}
	defer close()
	parse := newParser(scanner)

	parse.Parse()

	return name, len(triangles), nil
}
