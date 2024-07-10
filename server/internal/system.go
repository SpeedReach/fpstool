package internal

import (
	"fmt"
	"time"
)

type AimBotSystem struct {
	Source          Source
	MouseController MouseController
	FigureDetector  FigureDetector
}

func (system AimBotSystem) Start() {
	go system.Source.Start()
	started := false
	var start time.Time
	processed := 0
	for true {
		getLatestItem := func(ch <-chan ScreenShot) ScreenShot {
			var latest ScreenShot
			for {
				select {
				case latest = <-ch:
				default:
					return latest
				}
			}
		}

		latest := getLatestItem(system.Source.GetStream())
		if latest == (ScreenShot{}) {
			continue
		}
		if started == false {
			start = time.Now()
			started = true
		}

		println(fmt.Sprintf("Detecting %d", latest.Index))
		detected := system.FigureDetector.Detect(latest.Image)
		processed += 1
		for _, d := range detected {
			println(d.Type, d.X, d.Y, d.Confidence)
		}

		if processed == 100 {
			break
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Processed %d frames in %s\n", processed, elapsed)
	fmt.Printf("FPS: %f\n", float64(processed)/elapsed.Seconds())

}
