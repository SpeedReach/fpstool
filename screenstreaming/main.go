package main

import (
	"encoding/binary"
	"github.com/kbinani/screenshot"
	"image"
	"image/png"
	"net"
)

func main() {
	// Calculate where to capture
	bounds := screenshot.GetDisplayBounds(0)
	squareSize := min(bounds.Dx(), bounds.Dy()) / 5
	top := bounds.Dy()/2 - squareSize/2
	left := bounds.Dx()/2 - squareSize/2
	rect := image.Rect(left, top, left+squareSize, top+squareSize)
	tcp(rect)
}

func tcp(rect image.Rectangle) {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if err = binary.Write(conn, binary.LittleEndian, int32(rect.Dx())); err != nil {
		panic(err)
	}
	if err = binary.Write(conn, binary.LittleEndian, int32(rect.Dy())); err != nil {
		panic(err)
	}
	println("Wrote size to server")
	//start capturing and sending images
	for true {
		img, err := screenshot.CaptureRect(rect)
		if err != nil {
			panic(err)
		}
		if err = png.Encode(conn, img); err != nil {
			// retry connection
			continue
		}
		var buf = [1]byte{}
		if _, err = conn.Read(buf[0:]); err != nil {
			panic(err)
		}
		if buf[0] == 'S' {
			println("Server received image")
		} else {
			panic("Server did not receive image: " + string(buf[0]))
		}
	}
}

func udp(rect image.Rectangle) {

	// Send the capture size to the server
	conn, err := net.Dial("udp", "localhost:12345")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if err = binary.Write(conn, binary.LittleEndian, int32(rect.Dx())); err != nil {
		panic(err)
	}
	if err = binary.Write(conn, binary.LittleEndian, int32(rect.Dy())); err != nil {
		panic(err)
	}
	println("Wrote size to server")
	//start capturing and sending images
	for true {
		img, err := screenshot.CaptureRect(rect)
		if err != nil {
			panic(err)
		}
		if err = png.Encode(conn, img); err != nil {
			panic(err)
		}
		var buf = [1]byte{}
		if _, err = conn.Read(buf[0:]); err != nil {
			panic(err)
		}
		if buf[0] == 'S' {
			println("Server received image")
		} else {
			panic("Server did not receive image: " + string(buf[0]))
		}
	}
}
