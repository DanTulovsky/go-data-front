package obj

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

type ScanLine interface {
	Parse(*bufio.Scanner) error
	IsAtEOF() bool
	IsBlank() bool
	IsComment() bool
	IsCommand(name string) bool
	ParamCount() int
	StringParam(index int) string
	FloatParam(index int) (float32, error)
	CoordReferenceParam(index int) (ScanCoordReference, error)
	GetComment() string
}

type ScanCoordReference struct {
	HasTexCoordIndex bool
	HasNormalIndex   bool
	VertexIndex      int
	TexCoordIndex    int
	NormalIndex      int
}

type scanLine struct {
	atEOF        bool
	lineBuffer   bytes.Buffer
	isComment    bool
	comment      string
	segments     []string
	segmentRegex *regexp.Regexp
}

func NewScanLine() ScanLine {
	regex, err := regexp.Compile("[\\s]+")
	if err != nil {
		panic(err)
	}
	return &scanLine{
		atEOF:        false,
		lineBuffer:   bytes.Buffer{},
		isComment:    false,
		comment:      "",
		segments:     nil,
		segmentRegex: regex,
	}
}

func (c *scanLine) Parse(scanner *bufio.Scanner) error {
	c.comment = ""
	c.segments = []string{}

	logicalLine, err := c.readLogicalLine(scanner)
	if err != nil {
		return err
	}

	c.isComment = strings.HasPrefix(logicalLine, "#")
	if c.isComment {
		c.comment = strings.TrimSpace(strings.TrimPrefix(logicalLine, "#"))
	} else {
		c.segments = c.segmentRegex.Split(logicalLine, -1)
	}
	return nil
}

func (c *scanLine) IsAtEOF() bool {
	return c.atEOF
}

func (c *scanLine) IsBlank() bool {
	return (c.comment == "") && (c.segments == nil)
}

func (c *scanLine) IsComment() bool {
	return c.isComment
}

func (c *scanLine) IsCommand(name string) bool {
	return (len(c.segments) > 0) && (c.segments[0] == name)
}

func (c *scanLine) ParamCount() int {
	count := len(c.segments) - 1
	if count < 0 {
		count = 0
	}
	return count
}

func (c *scanLine) StringParam(index int) string {
	return c.segments[index+1]
}

func (c *scanLine) FloatParam(index int) (float32, error) {
	value, err := strconv.ParseFloat(c.StringParam(index), 32)
	if err != nil {
		return 0.0, err
	}
	return float32(value), nil
}

func (c *scanLine) CoordReferenceParam(index int) (ScanCoordReference, error) {
	var err error
	result := ScanCoordReference{}
	references := strings.Split(c.StringParam(index), "/")

	result.VertexIndex, err = strconv.Atoi(references[0])
	if err != nil {
		return ScanCoordReference{}, err
	}

	result.HasTexCoordIndex = len(references) > 1 && (references[1] != "")
	if result.HasTexCoordIndex {
		result.TexCoordIndex, err = strconv.Atoi(references[1])
		if err != nil {
			return ScanCoordReference{}, err
		}
	}

	result.HasNormalIndex = len(references) > 2 && (references[2] != "")
	if result.HasNormalIndex {
		result.NormalIndex, err = strconv.Atoi(references[2])
		if err != nil {
			return ScanCoordReference{}, err
		}
	}

	return result, nil
}

func (c *scanLine) GetComment() string {
	return strings.TrimSpace(c.lineBuffer.String()[1:])
}

func (c *scanLine) readLogicalLine(scanner *bufio.Scanner) (string, error) {
	c.lineBuffer.Reset()
	for scanner.Scan() {
		if scanner.Err() != nil {
			return "", scanner.Err()
		}
		line := scanner.Text()
		if !strings.HasSuffix(line, "\\") {
			c.lineBuffer.WriteString(line)
			return strings.TrimSpace(c.lineBuffer.String()), nil
		} else {
			content := strings.TrimSuffix(line, "\\")
			c.lineBuffer.WriteString(content)
		}
	}
	c.atEOF = true
	return strings.TrimSpace(c.lineBuffer.String()), nil
}
