package internal

import (
	"encoding/binary"
	"fmt"
	"image/png"
	"net"
	"time"
)

type TcpSource struct {
	started  bool
	pipeline chan ScreenShot
}

func (source TcpSource) GetStream() <-chan ScreenShot {
	return source.pipeline
}

func NewTcpSource() TcpSource {
	return TcpSource{
		started:  false,
		pipeline: make(chan ScreenShot, 10),
	}
}

func (source TcpSource) Start() {
	if source.started {
		return
	}
	source.started = true
	// Start a TCP server on port 8080
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("Server started on %s", listener.Addr().String())
	for {
		// Accept a new client connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the client connection in a new goroutine
		source.handleConnection(conn)
	}
}

func (source TcpSource) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr().String())
	var width, height int32
	if err := binary.Read(conn, binary.LittleEndian, &width); err != nil {
		panic(err)
	}
	if err := binary.Read(conn, binary.LittleEndian, &height); err != nil {
		panic(err)
	}
	fmt.Println("Received size:", width, height)
	//imgData := make([]byte, width*height*4)
	i := 0
	for {
		taken := time.Now()
		img, err := png.Decode(conn)
		if err != nil {
			panic(err)
		}
		i++
		source.pipeline <- ScreenShot{
			Taken: taken,
			Image: img,
			Index: i,
		}
		//println("piped image", i)
		if _, err := conn.Write([]byte("S")); err != nil {
			panic(err)
		}
	}
}
