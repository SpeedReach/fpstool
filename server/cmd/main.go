package main

import (
	"github.com/SpeedReach/fpstool/server/internal"
)

func main() {
	println("Starting")
	system := internal.AimBotSystem{
		Source: internal.NewTcpSource(),
		//MouseController: internal.NewSerialMouseController(),
		MouseReader:    internal.NewWindowsMouseReader(),
		FigureDetector: internal.NewYoloV5Detection(),
	}
	system.Start()
}
