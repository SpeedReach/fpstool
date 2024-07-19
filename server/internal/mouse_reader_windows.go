package internal

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procRegisterRawInput = user32.NewProc("RegisterRawInputDevices")
	procGetRawInputData  = user32.NewProc("GetRawInputData")
	procDefWindowProc    = user32.NewProc("DefWindowProcW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	procCreateWindowEx   = user32.NewProc("CreateWindowExW")
	procRegisterClass    = user32.NewProc("RegisterClassW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
)

const (
	RID_INPUTSINK = 0x00000100
	RIM_TYPEMOUSE = 0
	RID_INPUT     = 0x10000003
	RID_HEADER    = 0x10000005
	WM_INPUT      = 0x00FF
	WM_CREATE     = 0x0001
	WM_DESTROY    = 0x0002
	CW_USEDEFAULT = 0x80000000
	SW_SHOW       = 5

	RI_MOUSE_LEFT_BUTTON_DOWN   = 0x0001
	RI_MOUSE_LEFT_BUTTON_UP     = 0x0002
	RI_MOUSE_RIGHT_BUTTON_DOWN  = 0x0004
	RI_MOUSE_RIGHT_BUTTON_UP    = 0x0008
	RI_MOUSE_MIDDLE_BUTTON_DOWN = 0x0010
	RI_MOUSE_MIDDLE_BUTTON_UP   = 0x0020
)

type RAWINPUTDEVICE struct {
	UsagePage uint16
	Usage     uint16
	Flags     uint32
	Target    uintptr
}

type RAWINPUTHEADER struct {
	Type   uint32
	Size   uint32
	Device uintptr
	WParam uintptr
}

type RAWMOUSE struct {
	Header RAWINPUTHEADER
	Data   struct {
		Flags       uint16
		ButtonFlags uint16
		ButtonData  uint16
		RawButtons  uint32
		LastX       int32
		LastY       int32
		ExtraInfo   uint32
	}
}

type WindowsMouseReader struct {
	channel chan MouseEvent
}

func (ms WindowsMouseReader) ReadEvent() <-chan MouseEvent {
	return ms.channel
}

func (rawMouse RAWMOUSE) String() string {
	return fmt.Sprintf("Flags: %d, ButtonFlags: %d, ButtonData: %d, RawButtons: %d, LastX: %d, LastY: %d, ExtraInfo: %d", rawMouse.Data.Flags, rawMouse.Data.ButtonFlags, rawMouse.Data.ButtonData, rawMouse.Data.RawButtons, rawMouse.Data.LastX, rawMouse.Data.LastY, rawMouse.Data.ExtraInfo)
}

type WNDCLASS struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
}

type MSG struct {
	HWND    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct {
		X int32
		Y int32
	}
}

func registerRawInputDevices(hwnd uintptr) error {
	rid := RAWINPUTDEVICE{
		UsagePage: 0x01, // HID_USAGE_PAGE_GENERIC
		Usage:     0x02, // HID_USAGE_GENERIC_MOUSE
		Flags:     RID_INPUTSINK,
		Target:    hwnd,
	}

	_, _, err := procRegisterRawInput.Call(
		uintptr(unsafe.Pointer(&rid)),
		1,
		unsafe.Sizeof(rid),
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return err
	}

	return nil
}

func getRawInputData(lParam uintptr) (*RAWMOUSE, error) {
	var header RAWINPUTHEADER
	size := uint32(unsafe.Sizeof(header))
	_, _, err := procGetRawInputData.Call(
		lParam,
		RID_HEADER,
		uintptr(unsafe.Pointer(&header)),
		uintptr(unsafe.Pointer(&size)),
		uintptr(unsafe.Sizeof(header)),
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return nil, err
	}

	rawInput := make([]byte, header.Size)
	size = header.Size
	_, _, err = procGetRawInputData.Call(
		lParam,
		RID_INPUT,
		uintptr(unsafe.Pointer(&rawInput[0])),
		uintptr(unsafe.Pointer(&size)),
		uintptr(unsafe.Sizeof(header)),
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		return nil, err
	}

	raw := (*RAWMOUSE)(unsafe.Pointer(&rawInput[0]))
	return raw, nil
}

func (ms WindowsMouseReader) wndProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_CREATE:
		err := registerRawInputDevices(uintptr(hwnd))
		if err != nil {
			fmt.Println("Failed to register raw input device:", err)
		}
	case WM_INPUT:
		raw, err := getRawInputData(lParam)
		if err != nil {
			panic(err)
		} else {
			event := MouseEvent{X: int16(raw.Data.LastX), Y: int16(raw.Data.LastY)}
			if raw.Data.ButtonData&RI_MOUSE_LEFT_BUTTON_DOWN != 0 {
				event.LeftButton = PressDown
			}
			if raw.Data.ButtonData&RI_MOUSE_LEFT_BUTTON_UP != 0 {
				event.LeftButton = Release
			}
			if raw.Data.ButtonData&RI_MOUSE_RIGHT_BUTTON_DOWN != 0 {
				event.RightButton = PressDown
			}
			if raw.Data.ButtonData&RI_MOUSE_RIGHT_BUTTON_UP != 0 {
				event.RightButton = Release
			}
			ms.channel <- event
		}
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
	default:
		ret, _, _ := procDefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}
	return 0
}

func NewWindowsMouseReader() MouseReader {
	channel := make(chan MouseEvent)
	reader := WindowsMouseReader{
		channel: channel,
	}

	return reader
}
func (ms WindowsMouseReader) Start() {

	instance, _, _ := procGetModuleHandle.Call(0)
	className := syscall.StringToUTF16Ptr("RawInputClass")

	wc := WNDCLASS{
		Style:         0,
		LpfnWndProc:   syscall.NewCallback(ms.wndProc),
		HInstance:     syscall.Handle(instance),
		LpszClassName: className,
	}

	_, _, err := procRegisterClass.Call(uintptr(unsafe.Pointer(&wc)))
	if err != nil && err.Error() != "The operation completed successfully." {
		panic(err)
	}

	hwnd, _, err := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Raw Input Example"))),
		0,
		CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT,
		0,
		0,
		instance,
		0,
	)

	if err != nil && err.Error() != "The operation completed successfully." {
		panic(err)
	}

	procShowWindow.Call(hwnd, SW_SHOW)

	var msg MSG
	for {
		ret, _, err := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		if int(ret) == -1 {
			fmt.Println("Error in GetMessage:", err)
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
