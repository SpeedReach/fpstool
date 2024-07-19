package internal

import (
	"fmt"
)

type MouseReader interface {
	ReadEvent() <-chan MouseEvent
	Start()
}

type MouseButtonAction int8

const (
	Remain MouseButtonAction = iota
	PressDown
	Release
)

type MouseEvent struct {
	X           int16
	Y           int16
	LeftButton  MouseButtonAction
	RightButton MouseButtonAction
}

func (e MouseEvent) String() string {
	return fmt.Sprintf("X: %d, Y: %d, LeftButton: %d, RightButton: %d", e.X, e.Y, e.LeftButton, e.RightButton)
}
