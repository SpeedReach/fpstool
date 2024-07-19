package internal

import (
	"github.com/go-vgo/robotgo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
	"os"
)

type EbitenMouseReader struct {
	channel                      chan MouseEvent
	screenCenterX, screenCenterY int
	state                        *GameState
}

type GameState struct {
	leftPressed      bool
	rightPressed     bool
	mouseInitialized bool
	gameCenterX      int
	gameCenterY      int
}

func (e EbitenMouseReader) Start() {
	robotgo.Move(e.screenCenterX, e.screenCenterY)
	err := ebiten.RunGame(e)
	if err != nil {
		os.Exit(0)
	}
}

func (e EbitenMouseReader) initializeMouse() {
	posX, posY := ebiten.CursorPosition()
	e.state.gameCenterX = posX
	e.state.gameCenterY = posY
}

func (e EbitenMouseReader) ReadEvent() <-chan MouseEvent {
	return e.channel
}

func (e EbitenMouseReader) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		println("Exiting")
		os.Exit(0)
	}
	if !e.state.mouseInitialized {
		e.state.mouseInitialized = true
		e.initializeMouse()
		return nil
	}
	leftPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	rightPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	event := MouseEvent{}
	if e.state.leftPressed && !leftPressed {
		event.LeftButton = Release
	} else if !e.state.leftPressed && leftPressed {
		event.LeftButton = PressDown
	} else {
		event.LeftButton = Remain
	}

	if e.state.rightPressed && !rightPressed {
		event.RightButton = Release
	} else if !e.state.rightPressed && rightPressed {
		event.RightButton = PressDown
	} else {
		event.RightButton = Remain
	}

	e.state.leftPressed = leftPressed
	e.state.rightPressed = rightPressed

	posX, posY := ebiten.CursorPosition()
	if posX != e.state.gameCenterX || posY != e.state.gameCenterY {
		event.X = int16(posX - e.state.gameCenterX)
		event.Y = int16(posY - e.state.gameCenterY)
		robotgo.Move(e.screenCenterX, e.screenCenterY)
	}
	if event == (MouseEvent{}) {
		return nil
	}
	e.channel <- event
	return nil
}

func (e EbitenMouseReader) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Press ESC to exit")
}

func (e EbitenMouseReader) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func NewEbitenMouseReader() MouseReader {
	width, height := getScreenSize()
	channel := make(chan MouseEvent)
	reader := EbitenMouseReader{
		screenCenterX: width / 2,
		screenCenterY: height / 2,
		channel:       channel,
		state:         &GameState{},
	}
	ebiten.SetTPS(math.MaxInt)
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Mouse Reader")
	return reader
}

func getScreenSize() (int, int) {
	return robotgo.GetScreenSize()
}
