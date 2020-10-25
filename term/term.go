/*
Package term bundles functions and constants to navigate the terminal
*/
package term

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	tioCgWinSz = 0x5413 // On OSX use 1074295912. Thanks zeebo
)

//WinSize struct holds Row and Col for a terminal
type WinSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

//GetWinsize returns a pointer to a Winsize-Struct
func GetWinsize() (*WinSize, error) {
	ws := new(WinSize)

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tioCgWinSz),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		return nil, errors.New("WinSizeError:" + fmt.Sprint(errno))
	}
	return ws, nil
}
