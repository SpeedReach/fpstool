package internal

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"net"
	"os"
)

type UdpSource struct {
	stream  chan ScreenShot
	started bool
}

func (s UdpSource) GetStream() <-chan ScreenShot {
	//TODO implement me
	panic("implement me")
}

func NewUdpSource() UdpSource {
	return UdpSource{
		stream:  make(chan ScreenShot, 10),
		started: false,
	}
}

func (s UdpSource) Start() {
	if s.started {
		return
	}
	s.started = true
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	println("Listening on :12345")
	var width, height int32
	if err := binary.Read(conn, binary.LittleEndian, &width); err != nil {
		panic(err)
	}
	if err := binary.Read(conn, binary.LittleEndian, &height); err != nil {
		panic(err)
	}
	println("Received size: ", width, height)
	size := width * height * 3
	i := 0
	for {
		imgData := make([]byte, size)
		if _, _, err := conn.ReadFromUDP(imgData); err != nil {
			panic(err)
		}
		img, err := png.Decode(bytes.NewReader(imgData))
		if err != nil {
			panic(err)
		}
		i++
		file, err := os.Create("build/img" + string(rune(i)) + ".png")
		if err != nil {
			panic(err)
		}
		if err := png.Encode(file, img); err != nil {
			panic(err)
		}
		println("Image received")
		if _, err = conn.Write([]byte("S")); err != nil {
			panic(err)
		}
	}
}
