package internal

import "image"

type FigureDetector interface {
	Detect(data image.Image) []Detected
}

type Detected struct {
	Type       DetectedType
	X, Y       int
	Confidence float32
}

type DetectedType int

const (
	DetectedTypeBody DetectedType = iota
	DetectedTypeHead
)
