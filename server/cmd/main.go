package main

import "github.com/SpeedReach/fpstool/server/internal"

func main() {
	system := internal.AimBotSystem{
		Source:          internal.NewTcpSource(),
		MouseController: nil,
		FigureDetector:  internal.NewYoloV5Detection(),
	}
	system.Start()
}
