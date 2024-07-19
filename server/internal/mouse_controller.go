package internal

import (
	"bytes"
	"encoding/binary"
	"github.com/tarm/serial"
)

type MouseController interface {
	ControlMouse(event MouseEvent) error
}

type serialMouseController struct {
	port *serial.Port
}

func NewSerialMouseController() MouseController {
	c := &serial.Config{Name: "COM8", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	return serialMouseController{port: s}
}

func (mc serialMouseController) ControlMouse(event MouseEvent) error {
	mc.send(toMessage(event))
	mc.waitForAck()
	return nil
}

type message struct {
	start byte // 255
	dx    int16
	dy    int16
	left  int8
	right int8
}

func toMessage(event MouseEvent) message {
	return message{start: 255, dx: event.X, dy: event.Y, left: int8(event.LeftButton), right: int8(event.RightButton)}
}

func (mc serialMouseController) send(mes message) {
	//println("Sending message")
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, mes)
	//fmt.Printf("binary: %v\n", buf.Bytes())
	_, err := mc.port.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
	//println("Sent")
}

func (mc serialMouseController) waitForAck() {

	recv := make([]byte, 1)
	_, err := mc.port.Read(recv)
	if err != nil {
		panic(err)
	}
}

func (mc serialMouseController) getDebugMsg() {
	recv := make([]byte, 31)
	_, err := mc.port.Read(recv)
	if err != nil {
		panic(err)
	}
	println(string(recv))
}
