package internal

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"
)

type AimBotSystem struct {
	Source          Source
	MouseController MouseController
	MouseReader     MouseReader
	FigureDetector  FigureDetector
}

func (system AimBotSystem) Start() {
	go system.Source.Start()
	go func() {
		if system.MouseController == nil {
			println("No mouse controller")
			return
		}
		for movement := range system.MouseReader.ReadEvent() {
			//println(movement.String())

			_ = system.MouseController.ControlMouse(movement)
		}
	}()

	go func() {
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

			//println(fmt.Sprintf("Detecting %d", latest.Index))
			detected := system.FigureDetector.Detect(latest.Image)
			processed += 1

			//rgba := toRGBA(latest.Image)
			for _, d := range detected {
				if d.Confidence > 0.4 {
					log.Printf("%d, %d, %d, %f", d.Type, d.X, d.Y, d.Confidence)
				}
			}
		}

		elapsed := time.Since(start)
		fmt.Printf("Processed %d frames in %s\n", processed, elapsed)
		fmt.Printf("FPS: %f\n", float64(processed)/elapsed.Seconds())
	}()

	system.MouseReader.Start()
}

func toRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	// Fill the image with the fill color
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgba.SetRGBA(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}
	return rgba
}
