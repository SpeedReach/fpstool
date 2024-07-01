package main

import "github.com/kbinani/screenshot"

func main() {
	bounds := screenshot.GetDisplayBounds(0)
	width := bounds.Size().X

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
	file, _ := os.Create(fileName)
	defer file.Close()
	png.Encode(file, img)

	fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
}
