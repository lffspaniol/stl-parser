package stl

// This file defines a parser for the STL ASCII format.

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var ErrEOF = errors.New("File is empty")
var ErrInvalidExpectedToken = errors.New(`"facet" or "endsolid" expected`)
var ErrInvalidSintax = errors.New("Invalid sintax")
var ErrInvalidFloat = errors.New("Invalid float")
var errInvalidToken = errors.New("Invalid token")

var expectedASCIIHeaderPrefix = []byte("solid ")

const (
	idNone  = 0
	idSolid = 1 << iota
	idFacet
	idNormal
	idOuter
	idLoop
	idVertex
	idEndloop
	idEndfacet
	idEndsolid
)

var identRegexps = map[int]*regexp.Regexp{
	idSolid:              regexp.MustCompile("^solid$"),
	idFacet:              regexp.MustCompile("^facet$"),
	idNormal:             regexp.MustCompile("^normal$"),
	idOuter:              regexp.MustCompile("^outer$"),
	idLoop:               regexp.MustCompile("^loop$"),
	idVertex:             regexp.MustCompile("^vertex$"),
	idEndloop:            regexp.MustCompile("^endloop$"),
	idEndfacet:           regexp.MustCompile("^endfacet$"),
	idEndsolid:           regexp.MustCompile("^endsolid$"),
	idFacet | idEndsolid: regexp.MustCompile(`^(facet|endsolid)$`),
}

var idents = map[int]string{
	idSolid:    "solid",
	idFacet:    "facet",
	idNormal:   "normal",
	idOuter:    "outer",
	idLoop:     "loop",
	idVertex:   "vertex",
	idEndloop:  "endloop",
	idEndfacet: "endfacet",
	idEndsolid: "endsolid",
}

func extractASCIIString(byteData []byte) string {
	i := 0
	for i < len(byteData) && byteData[i] < byte(128) && byteData[i] != byte(0) {
		i++
	}
	return string(byteData[0:i])
}

type parser struct {
	line             int
	errors           error
	currentWord      string
	currentLine      []byte
	eof              bool
	lineScanner      *bufio.Scanner
	wordScanner      *bufio.Scanner
	HeaderError      bool
	TrianglesSkipped bool
	ErrorText        string
}

func (p *parser) addError(msg error) {
	p.errors = errors.Join(fmt.Errorf("%d: %w", p.line, msg))
}

func (p *parser) Parse() bool {
	if p.eof {
		p.HeaderError = true
		p.addError(ErrEOF)
	} else {
		p.HeaderError = !p.parseASCIIHeaderLine()
	TriangleLoop:
		for !p.eof && !p.isCurrentTokenIdent(idEndsolid) {
			if !p.isCurrentTokenIdent(idFacet) {
				p.addError(ErrInvalidExpectedToken)
				switch p.skipToToken(idFacet | idEndsolid) {
				case idEndsolid, idNone:
					break TriangleLoop
				}
			}

			var t Triangle
			if p.parseFacet(&t) {
				triangles = append(triangles, t)
			} else {
				p.TrianglesSkipped = true
				p.skipToToken(idFacet | idEndsolid)
			}
		}
	}

	success := !p.HeaderError && !p.TrianglesSkipped && p.consumeToken(idEndsolid)
	return success
}

func (p *parser) parseASCIIHeaderLine() bool {
	var success bool
	if p.eof {
		p.addError(ErrEOF)
		success = false
	} else {
		if !bytes.HasPrefix(p.currentLine, expectedASCIIHeaderPrefix) {
			p.addError(ErrInvalidSintax)
			success = false
		} else {
			name = extractASCIIString(p.currentLine[len(expectedASCIIHeaderPrefix):])
			success = true
		}
	}
	p.nextLine()
	return success
}

func (p *parser) parseFacet(t *Triangle) bool {
	return p.consumeToken(idFacet) &&
		p.consumeToken(idNormal) && p.parsePoint(&t.Normal) &&
		p.consumeToken(idOuter) && p.consumeToken(idLoop) &&
		p.consumeToken(idVertex) && p.parsePoint(&t.Vertices[0]) &&
		p.consumeToken(idVertex) && p.parsePoint(&t.Vertices[1]) &&
		p.consumeToken(idVertex) && p.parsePoint(&t.Vertices[2]) &&
		p.consumeToken(idEndloop) &&
		p.consumeToken(idEndfacet)
}

func (p *parser) parsePoint(pt *[3]float32) bool {
	return p.parseFloat32(&pt[0]) &&
		p.parseFloat32(&pt[1]) &&
		p.parseFloat32(&pt[2])
}

func (p *parser) parseFloat32(f *float32) bool {
	if p.eof {
		return false
	}
	f64, err := strconv.ParseFloat(p.currentWord, 32)
	if err != nil {
		p.addError(ErrInvalidFloat)
		return false
	}

	*f = float32(f64)
	p.nextWord()
	return true
}

func (p *parser) isCurrentTokenIdent(ident int) bool {
	re := identRegexps[ident]
	return re.MatchString(p.currentWord)
}

func (p *parser) skipToToken(ident int) int {
	re := identRegexps[ident]
	for { // terminates when no more next words are there, or ident has been found
		if re.MatchString(p.currentWord) {
			if ident == (idFacet | idEndsolid) {
				if identRegexps[idFacet].MatchString(p.currentWord) {
					return idFacet
				}
				return idEndsolid
			}
			return ident
		}
		if !p.nextWord() {
			return idNone
		}
	}
}

func (p *parser) consumeToken(ident int) bool {
	re := identRegexps[ident]
	if !re.MatchString(p.currentWord) {
		ident := idents[ident]
		p.addError(fmt.Errorf(`%w: %s`, errInvalidToken, ident))
		return false
	}

	p.nextWord()
	return true
}

func (p *parser) nextWord() bool {
	if p.eof {
		return false
	}
	// Try to advance word scanner
	if p.wordScanner.Scan() {
		p.currentWord = p.wordScanner.Text()
		return true
	}
	if p.wordScanner.Err() == nil { // line has ended
		return p.nextLine()
	}
	p.addError(p.wordScanner.Err())
	p.currentLine = nil
	p.currentWord = ""
	p.eof = true
	return false
}

func (p *parser) nextLine() bool {
	if p.lineScanner.Scan() {
		p.currentLine = p.lineScanner.Bytes()
		p.line++
		p.wordScanner = bufio.NewScanner(bytes.NewReader(p.currentLine))
		p.wordScanner.Split(bufio.ScanWords)
		return p.nextWord()
	}

	if p.lineScanner.Err() != nil {
		p.addError(p.lineScanner.Err())
	}
	p.currentLine = nil
	p.currentWord = ""
	p.eof = true
	return false
}

func newParser(lineScanner *bufio.Scanner) *parser {
	var p parser
	p.eof = false
	p.lineScanner = lineScanner
	p.nextLine()
	return &p
}
