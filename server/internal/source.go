package internal

import (
	"image"
	"time"
)

type Source interface {
	Start()
	GetStream() <-chan ScreenShot
}

type ScreenShot struct {
	Taken time.Time
	Image image.Image
	Index int
}
